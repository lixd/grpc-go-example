
# hello world demo
gen-hello-world:
	protoc --proto_path=./helloworld/helloworld \
   --go_out=./helloworld/helloworld --go_opt=paths=source_relative \
   --go-grpc_out=./helloworld/helloworld --go-grpc_opt=paths=source_relative \
   ./helloworld/helloworld/hello_world.proto

# echo demo
gen-echo:
	protoc --proto_path=./features/proto \
   --go_out=./features/proto --go_opt=paths=source_relative \
   --go-grpc_out=./features/proto  --go-grpc_opt=paths=source_relative \
   ./features/proto/echo/echo.proto

# grpc-gateway demo
gen-gw:
	protoc --proto_path=./features/proto \
   --go_out=./features/proto --go_opt=paths=source_relative \
   --go-grpc_out=./features/proto --go-grpc_opt=paths=source_relative \
   --grpc-gateway_out=./features/proto --grpc-gateway_opt=paths=source_relative \
   ./features/proto/gateway/gateway.proto
# proto import demo
gen-imp:
	protoc --proto_path=. --go_out=. ./protobuf/import/*.proto
