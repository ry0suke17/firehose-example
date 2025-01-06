package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	fetchpb "github.com/ry0suke17/firehose-example/proto/sf/firehose/v2"
	typepb "github.com/ry0suke17/firehose-example/proto/sf/solana/type/v1"
)

// blocksPerFile represents the number of blocks contained in each merged .dbin.zst file.
const blocksPerFile = 100

// latestSafeBlock is a simple global variable to track the largest block number
// that we consider "safe to remove" from the merged blocks.
var latestSafeBlock uint64

func main() {
	//-------------------------------------------------------
	// 1) Flags for command-line arguments
	//-------------------------------------------------------
	serverURL := flag.String("url", "localhost:10015", "gRPC server URL")
	startBlock := flag.Uint64("block", 309000000, "Block number to start fetching from")

	// A new flag for the firehose data directory (default = ./firehose-data)
	dataDir := flag.String("data-dir", "./firehose-data", "Root directory for Firehose data")

	// Get interval
	getInterval := flag.Duration("get-interval", 500*time.Millisecond, "Interval for get blocks")

	// Cleanup interval
	cleanupInterval := flag.Duration("cleanup-interval", 1*time.Minute, "Interval for cleanupOldMergedBlocks")
	flag.Parse()

	//-------------------------------------------------------
	// 2) Build paths for merged-blocks and one-blocks
	//-------------------------------------------------------
	mergedBlocksDir := filepath.Join(*dataDir, "storage", "merged-blocks")
	oneBlocksDir := filepath.Join(*dataDir, "storage", "one-blocks")

	//-------------------------------------------------------
	// 3) Signal handling
	//-------------------------------------------------------
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	done := make(chan struct{})

	go func() {
		<-sigCh
		fmt.Println("Received signal, shutting down gracefully...")
		close(done)
	}()

	//-------------------------------------------------------
	// 4) gRPC connection
	//-------------------------------------------------------
	maxMsgSize := 10 * 1024 * 1024 // 10MB
	conn, err := grpc.Dial(
		*serverURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := fetchpb.NewFetchClient(conn)

	//-------------------------------------------------------
	// 5) Initialization
	//-------------------------------------------------------
	latestSafeBlock = *startBlock

	//-------------------------------------------------------
	// 6) Start a goroutine for cleanup
	//-------------------------------------------------------
	go func() {
		ticker := time.NewTicker(*cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				log.Println("[CLEANUP] Received done, stop cleaner goroutine...")
				return
			case <-ticker.C:
				// 1) Collect the block numbers in the one-blocks directory
				oneBlockNumbers, err := collectOneBlockNumbers(oneBlocksDir)
				if err != nil {
					log.Printf("[CLEANUP] failed to collect one-block numbers: %v\n", err)
					continue
				}

				// 2) Perform cleanup
				if err := cleanupOldMergedBlocks(mergedBlocksDir, latestSafeBlock, oneBlockNumbers); err != nil {
					log.Printf("[CLEANUP] error: %v\n", err)
				}
			}
		}
	}()

	//-------------------------------------------------------
	// 7) Main loop (block fetching)
	//-------------------------------------------------------
	currentBlock := *startBlock

	for {
		select {
		case <-done:
			fmt.Println("Stopping block fetch loop...")
			return
		default:
			// no signal yet
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		req := &fetchpb.SingleBlockRequest{
			Reference: &fetchpb.SingleBlockRequest_BlockNumber_{
				BlockNumber: &fetchpb.SingleBlockRequest_BlockNumber{
					Num: currentBlock,
				},
			},
		}
		resp, rpcErr := client.Block(ctx, req)
		cancel()

		if rpcErr != nil {
			// Check if it's a gRPC error with a specific code
			st, ok := status.FromError(rpcErr)
			if ok && st.Code() == codes.NotFound {
				// If the server responds with NotFound, skip this block
				log.Printf("[Block %d] NotFound. Skipping block...", currentBlock)
				currentBlock++
				time.Sleep(*getInterval)
				continue
			} else {
				// For other errors, just log and continue
				log.Printf("[Block %d] Block RPC failed: %v", currentBlock, rpcErr)
				time.Sleep(*getInterval)
				continue
			}
		} else {
			log.Printf("[Block %d] Got SingleBlockResponse\n", currentBlock)
			if resp.Block != nil {
				var myBlock typepb.Block
				if err := anypb.UnmarshalTo(resp.Block, &myBlock, proto.UnmarshalOptions{}); err != nil {
					log.Printf("failed to unmarshal block: %v\n", err)
				} else {
					if myBlock.Slot < currentBlock {
						log.Printf("[Block %d] Unexpected block slot = %d, skipping...",
							currentBlock, myBlock.Slot)
						time.Sleep(2 * time.Second)
						continue
					} else {
						log.Printf("[Block %d] OK, tx length = %d\n",
							currentBlock, len(myBlock.Transactions))
					}
				}
			}
			// Additional handling for resp.Metadata if needed
		}

		latestSafeBlock = currentBlock

		currentBlock++
		time.Sleep(*getInterval)
	}
}

// collectOneBlockNumbers scans the oneBlocksDir for any .dbin.zst files
// and extracts the block number from the filename, returning a set of
// block numbers that are still present in one-blocks directory.
func collectOneBlockNumbers(oneBlocksDir string) (map[uint64]bool, error) {
	blockNums := make(map[uint64]bool)

	err := filepath.WalkDir(oneBlocksDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".dbin.zst") {
			return nil
		}

		// Example filename: 0309000300-some-hash-309000299-default.dbin.zst
		parts := strings.Split(strings.TrimSuffix(d.Name(), ".dbin.zst"), "-")
		if len(parts) < 1 {
			return nil
		}

		blockStr := parts[0]

		num, parseErr := strconv.ParseUint(blockStr, 10, 64)
		if parseErr != nil {
			return nil
		}
		blockNums[num] = true
		return nil
	})

	if err != nil {
		return nil, err
	}
	return blockNums, nil
}

// cleanupOldMergedBlocks checks mergedBlocksDir for .dbin.zst files and deletes them
// if endBlock < safeBlock, unless that block number is found in oneBlockNumbers set.
func cleanupOldMergedBlocks(dir string, safeBlock uint64, oneBlockNumbers map[uint64]bool) error {
	log.Printf("[CLEANUP] start. safeBlock=%d", safeBlock)

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".dbin.zst") {
			return nil
		}

		filenameNoExt := strings.TrimSuffix(d.Name(), ".dbin.zst")
		filenameNoExt = strings.TrimSuffix(filenameNoExt, ".dbin")

		startBlock, parseErr := strconv.ParseUint(filenameNoExt, 10, 64)
		if parseErr != nil {
			return nil
		}

		endBlock := startBlock + blocksPerFile - 1

		if oneBlockNumbers[startBlock] {
			log.Printf("[CLEANUP] skip removing %s because block %d is still in one-blocks",
				path, startBlock)
			return nil
		}

		if endBlock < safeBlock {
			log.Printf("[CLEANUP] removing file: %s (blocks %d~%d < %d)",
				path, startBlock, endBlock, safeBlock)
			if delErr := os.Remove(path); delErr != nil {
				log.Printf("[CLEANUP] failed to remove %s: %v", path, delErr)
			}
		}

		log.Printf("[CLEANUP] no action for %s (blocks %d~%d >= %d)",
			path, startBlock, endBlock, safeBlock)

		return nil
	})
}
