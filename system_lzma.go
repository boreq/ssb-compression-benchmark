package ssb_compression_benchmark

import (
	"bytes"
	"github.com/boreq/errors"
	"github.com/ulikunitz/xz/lzma"
	"io"
)

type SystemLZMA struct {
	buf *bytes.Buffer
}

func NewSystemLZMA() (*SystemLZMA, error) {
	return &SystemLZMA{
		buf: bytes.NewBuffer(nil),
	}, nil
}

func (s SystemLZMA) Compress(in []byte) ([]byte, error) {
	s.buf.Reset()

	writer, err := lzma.NewWriter(s.buf)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a writer")
	}
	n, err := writer.Write(in)
	if err != nil {
		return nil, errors.Wrap(err, "error writing")
	}
	if n != len(in) {
		return nil, errors.New("invalid length")
	}

	if err := writer.Close(); err != nil {
		return nil, errors.Wrap(err, "error closing")
	}

	return s.buf.Bytes(), nil
}

func (s SystemLZMA) Decompress(in []byte) ([]byte, error) {
	s.buf.Reset()

	reader, err := lzma.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, errors.Wrap(err, "error creating a reader")
	}

	if _, err := io.Copy(s.buf, reader); err != nil {
		return nil, errors.Wrap(err, "error copying")
	}

	return s.buf.Bytes(), nil
}

type SystemLZMA2 struct {
	buf *bytes.Buffer
}

func NewSystemLZMA2() (*SystemLZMA2, error) {
	return &SystemLZMA2{
		buf: bytes.NewBuffer(nil),
	}, nil
}

func (s SystemLZMA2) Compress(in []byte) ([]byte, error) {
	s.buf.Reset()

	writer, err := lzma.NewWriter2(s.buf)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a writer")
	}
	n, err := writer.Write(in)
	if err != nil {
		return nil, errors.Wrap(err, "error writing")
	}
	if n != len(in) {
		return nil, errors.New("invalid length")
	}

	if err := writer.Close(); err != nil {
		return nil, errors.Wrap(err, "error closing")
	}

	return s.buf.Bytes(), nil
}

func (s SystemLZMA2) Decompress(in []byte) ([]byte, error) {
	s.buf.Reset()

	reader, err := lzma.NewReader2(bytes.NewReader(in))
	if err != nil {
		return nil, errors.Wrap(err, "error creating a reader")
	}

	if _, err := io.Copy(s.buf, reader); err != nil {
		return nil, errors.Wrap(err, "error copying")
	}

	return s.buf.Bytes(), nil
}
