proto-gen-gateway:
	protoc -I ./proto \
  --go_out ./pb --go_opt paths=source_relative \
  --go-grpc_out ./pb --go-grpc_opt paths=source_relative \
  --openapiv2_out=docs/apis/swagger \
  --grpc-gateway_out ./pb --grpc-gateway_opt paths=source_relative \
  proto/*.proto
gen-no-swagger:
  protoc -I ./proto \
    --go_out ./pb --go_opt paths=source_relative \
    --go-grpc_out ./pb --go-grpc_opt paths=source_relative \
    --grpc-gateway_out ./pb --grpc-gateway_opt paths=source_relative \
    proto/*.proto
run:
  go run components/story_service/main.go

run-client:
  go run ~/89minutes//components/story_service/client/main.go upload -a localhost:8080 -d test/blog/

  