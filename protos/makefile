GEN_DIR = gen

.PHONY: generate
generate:
	@echo "Generating protobuf code..."
	protoc --proto_path=.  --go_out=./$(GEN_DIR)  --go_opt=paths=source_relative  --go-grpc_out=./$(GEN_DIR)  --go-grpc_opt=paths=source_relative proto/hotel.proto
	@echo "Protobuf code generated in $(GEN_DIR)"