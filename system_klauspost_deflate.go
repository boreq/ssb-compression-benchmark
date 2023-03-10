package ssb_compression_benchmark

import (
	"bytes"
	"github.com/boreq/errors"
	"github.com/klauspost/compress/flate"
	"io"
)

type KlauspostDeflateLevel int

const (
	KlauspostDeflateLevelBestCompression    KlauspostDeflateLevel = flate.BestCompression
	KlauspostDeflateLevelDefaultCompression KlauspostDeflateLevel = flate.DefaultCompression
	KlauspostDeflateLevelBestSpeed          KlauspostDeflateLevel = flate.BestSpeed
)

type SystemKlauspostDeflate struct {
	level KlauspostDeflateLevel
}

func NewSystemKlauspostDeflate(level KlauspostDeflateLevel) (*SystemKlauspostDeflate, error) {
	return &SystemKlauspostDeflate{level: level}, nil
}

func (s SystemKlauspostDeflate) Compress(in []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	writer, err := flate.NewWriter(buf, int(s.level))
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

	return buf.Bytes(), nil
}

func (s SystemKlauspostDeflate) Decompress(in []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	reader := flate.NewReader(bytes.NewReader(in))

	if _, err := io.Copy(buf, reader); err != nil {
		return nil, errors.Wrap(err, "error copying")
	}

	return buf.Bytes(), nil
}
