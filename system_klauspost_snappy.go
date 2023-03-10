package ssb_compression_benchmark

import (
	"github.com/klauspost/compress/snappy"
)

type SystemKlauspostSnappy struct {
}

func NewSystemKlauspostSnappy() (*SystemKlauspostSnappy, error) {
	return &SystemKlauspostSnappy{}, nil
}

func (s *SystemKlauspostSnappy) Compress(in []byte) ([]byte, error) {
	return snappy.Encode(nil, in), nil
}

func (s *SystemKlauspostSnappy) Decompress(in []byte) ([]byte, error) {
	return snappy.Decode(nil, in)
}
