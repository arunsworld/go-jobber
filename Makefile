build:
	protoc --go_out=plugins=grpc:. api.proto