package ssb_compression_benchmark

import (
	"bytes"
	"github.com/boreq/errors"
	"github.com/klauspost/compress/s2"
)

type SystemKlauspostS2 struct {
	writer *s2.Writer
	reader *s2.Reader
	buf    *bytes.Buffer
}

func NewSystemKlauspostS2(option ...s2.WriterOption) (*SystemKlauspostS2, error) {
	return &SystemKlauspostS2{
		writer: s2.NewWriter(nil, option...),
		buf:    bytes.NewBuffer(nil),
	}, nil
}

func (s *SystemKlauspostS2) Compress(in []byte) ([]byte, error) {
	s.buf.Reset()
	s.writer.Reset(s.buf)

	n, err := s.writer.Write(in)
	if err != nil {
		return nil, errors.Wrap(err, "error writing")
	}

	if n != len(in) {
		return nil, errors.New("invalid length")
	}

	if err := s.writer.Close(); err != nil {
		return nil, errors.Wrap(err, "error closing")
	}

	return s.buf.Bytes(), nil
}

func (s *SystemKlauspostS2) Decompress(in []byte) ([]byte, error) {
	return s2.Decode(nil, in)
}
