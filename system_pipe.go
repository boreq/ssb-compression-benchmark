package ssb_compression_benchmark

import "github.com/boreq/errors"

type SystemPipe struct {
	systems []CompressionSystem
}

func NewSystemPipe(systems ...CompressionSystem) (*SystemPipe, error) {
	return &SystemPipe{systems: systems}, nil
}

func (s SystemPipe) Compress(in []byte) ([]byte, error) {
	for _, system := range s.systems {
		var err error
		in, err = system.Compress(in)
		if err != nil {
			return nil, errors.Wrap(err, "error compressing using provided system")
		}
	}
	return in, nil
}

func (s SystemPipe) Decompress(in []byte) ([]byte, error) {
	for i := len(s.systems) - 1; i >= 0; i-- {
		var err error
		in, err = s.systems[i].Decompress(in)
		if err != nil {
			return nil, errors.Wrap(err, "error compressing using provided system")
		}
	}
	return in, nil
}
