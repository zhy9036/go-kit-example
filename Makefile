compile-proto:
	protoc -I . --go_out=. --go-grpc_out=. ./protos/addsvc.proto