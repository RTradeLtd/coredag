# rebuild protocol buffers
.PHONY: proto
proto:
	protoc -I pb pb/pb.proto --go_out=plugins=grpc:pb