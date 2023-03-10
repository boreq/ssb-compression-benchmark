package ssb_compression_benchmark

import "github.com/klauspost/compress/snappy"

type SystemGolangSnappy struct {
}

func NewSystemGolangSnappy() (*SystemGolangSnappy, error) {
	return &SystemGolangSnappy{}, nil
}

func (s SystemGolangSnappy) Compress(in []byte) ([]byte, error) {
	return snappy.Encode(nil, in), nil
}

func (s SystemGolangSnappy) Decompress(in []byte) ([]byte, error) {
	return snappy.Decode(nil, in)
}
