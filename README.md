# firehose-example

Sample of getting block information with StreamingFast's Firehose.

## Usage

Start firehose grpc server.

```sh
make firehose-grpc
```

Activate firehose reader, merger, and relayer.

```sh
make firehose-grpc
```

Start the process of retrieve blocks. (and remove unnecessary 100-blocks)

```sh
make get-block
```