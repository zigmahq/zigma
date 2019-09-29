package version

//go:generate go get github.com/gogo/protobuf/protoc-gen-gofast
//go:generate protoc --gofast_out=. version.proto
//go:generate go run ../cmd/zsigner sign
