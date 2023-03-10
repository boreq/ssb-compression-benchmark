package ssb_compression_benchmark

import (
	"github.com/boreq/errors"
	"github.com/klauspost/compress/zstd"
)

type SystemKlauspostZSTD struct {
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func NewSystemKlauspostZSTD(level zstd.EncoderLevel) (*SystemKlauspostZSTD, error) {
	encoder, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(level))
	if err != nil {
		return nil, errors.Wrap(err, "error creating a writer")
	}

	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a reader")
	}

	return &SystemKlauspostZSTD{
		encoder: encoder,
		decoder: decoder,
	}, nil
}

func (s *SystemKlauspostZSTD) Compress(in []byte) ([]byte, error) {
	return s.encoder.EncodeAll(in, nil), nil
}

func (s *SystemKlauspostZSTD) Decompress(in []byte) ([]byte, error) {
	return s.decoder.DecodeAll(in, nil)
}
