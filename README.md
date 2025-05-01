
for generating protobuf for producer
```sh
protoc \
  --proto_path=. \
  --go_out=producers/grpc-producer/internal/proto \
  --go-grpc_out=producers/grpc-producer/internal/proto \
  $(find proto -name "*.proto")
  ```