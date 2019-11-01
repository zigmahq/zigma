package dht

import "github.com/multiformats/go-multihash"

// Hashable type
type Hashable interface {
	Data() []byte
	Hash() []byte
}

type hashable struct {
	data []byte
	hash []byte
}

func (h *hashable) Data() []byte {
	return h.data
}

func (h *hashable) Hash() []byte {
	return h.hash
}

// Bytes returns data in Hashable type
func Bytes(data []byte) Hashable {
	d := &hashable{
		data: data,
	}
	h, err := multihash.Sum(d.data, h, -1)
	if err != nil {
		panic(err)
	}
	d.hash = h
	return d
}

// String returns data in Hashable type
func String(data string) Hashable {
	return Bytes([]byte(data))
}
