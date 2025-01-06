proto:
	buf generate

DATA_DIR=firehose-data-test2

firehose-grpc:
	 firecore start -c firehose.yaml -d $(DATA_DIR)

firehose-block:
	firecore start -c reader-merger-relayer.yaml -d $(DATA_DIR)

get-block:
	go run cmd/main.go \
		--url=localhost:10015 \
		--block=309000000 \
		--data-dir=$(DATA_DIR) \
		--get-interval=100ms \
		--cleanup-interval=10s