package ssb_compression_benchmark

type CompressionSystem interface {
	Compress(in []byte) ([]byte, error)
	Decompress(in []byte) ([]byte, error)
}
