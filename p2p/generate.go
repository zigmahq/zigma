package p2p

//go:generate go get github.com/gogo/protobuf/protoc-gen-gofast
//go:generate protoc --gofast_out=plugins=grpc:. message.proto
