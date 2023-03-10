package ssb_compression_benchmark

import (
	"bytes"
	"github.com/andybalholm/brotli"
	"github.com/boreq/errors"
	"io"
)

type SystemBrotli struct {
	writer *brotli.Writer
}

func NewSystemBrotli(level int) (*SystemBrotli, error) {
	return &SystemBrotli{
		writer: brotli.NewWriterLevel(nil, level),
	}, nil
}

func (s *SystemBrotli) Compress(in []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	s.writer.Reset(buf)

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

	return buf.Bytes(), nil
}

func (s *SystemBrotli) Decompress(in []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	reader := brotli.NewReader(bytes.NewReader(in))

	if _, err := io.Copy(buf, reader); err != nil {
		return nil, errors.Wrap(err, "error copying")
	}

	return buf.Bytes(), nil
}
