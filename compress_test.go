package ssb_compression_benchmark_test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/andybalholm/brotli"
	ssb_compression "github.com/boreq/ssb-compression-benchmark"
	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type TestCompressionSystem struct {
	Name   string
	Create func() (ssb_compression.CompressionSystem, error)
}

const (
	BENCH_SNAPPY    = "BENCH_SNAPPY"
	BENCH_S2        = "BENCH_S2"
	BENCH_S2_SNAPPY = "BENCH_S2_SNAPPY"
	BENCH_DEFLATE   = "BENCH_DEFLATE"
	BENCH_LZMA      = "BENCH_LZMA"
	BENCH_ZSTD      = "BENCH_ZSTD"
	BENCH_BROTLI    = "BENCH_BROTLI"
)

type MessageKinds struct {
	FeedMessages                int
	CreateHistoryStreamRequests int
	OtherRequests               int
	EbtNotes                    int
	Other                       int
}

const (
	sampleSize              = 0.025
	steadySamplesMultiplier = 20
)

func TestMessageSources(t *testing.T) {
	messageSources := []struct {
		Name       string
		Path       string
		SampleSize float64
	}{
		{
			Name:       "many_feed_messages",
			Path:       "initial.txt",
			SampleSize: sampleSize,
		},
		{
			Name:       "few_feed_messages",
			Path:       "steady.txt",
			SampleSize: steadySamplesMultiplier * sampleSize,
		},
	}

	for _, messageSource := range messageSources {
		t.Run(messageSource.Name, func(t *testing.T) {
			messages, err := Load(testdataFilepath(messageSource.Path))
			require.NoError(t, err)

			kinds := MessageKinds{}

			for _, msg := range messages {
				p := msg.Payload

				if bytes.Contains(p, []byte(`"signature"`)) {
					kinds.FeedMessages++
					continue
				}

				if bytes.Contains(p, []byte(`"name":["createHistoryStream"]`)) {
					kinds.CreateHistoryStreamRequests++
					continue
				}

				if bytes.Contains(p, []byte(`{"name":[`)) {
					kinds.OtherRequests++
					continue
				}

				if bytes.Contains(p, []byte(`.ed25519":`)) {
					kinds.EbtNotes++
					continue
				}

				kinds.Other++
			}

			t.Log(fmt.Sprintf("%+v", kinds))

			messages = messages[:int(float64(len(messages))*messageSource.SampleSize)]

			t.Log("source:", messageSource.Name, "number of messages:", len(messages))
		})
	}
}

func BenchmarkLines(b *testing.B) {
	batches := []int{1, 10, 100}

	messageSources := []struct {
		Name       string
		Path       string
		SampleSize float64
	}{
		{
			Name:       "many_feed_messages",
			Path:       "initial.txt",
			SampleSize: sampleSize,
		},
		{
			Name:       "few_feed_messages",
			Path:       "steady.txt",
			SampleSize: steadySamplesMultiplier * sampleSize,
		},
	}

	var compressionSystems []TestCompressionSystem

	benchAll := true
	for _, trigger := range []string{
		BENCH_SNAPPY,
		BENCH_S2,
		BENCH_S2_SNAPPY,
		BENCH_DEFLATE,
		BENCH_LZMA,
		BENCH_ZSTD,
		BENCH_BROTLI,
	} {
		if os.Getenv(trigger) != "" {
			benchAll = false
		}
	}

	if benchAll || os.Getenv(BENCH_SNAPPY) != "" {
		compressionSystems = append(compressionSystems, []TestCompressionSystem{
			{
				Name: "snappy",
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemKlauspostSnappy()
				},
			},
		}...)
	}

	if benchAll || os.Getenv(BENCH_S2) != "" {
		compressionSystems = append(compressionSystems, []TestCompressionSystem{
			{
				Name: "s2_default",
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemKlauspostS2()
				},
			},
			{
				Name: "s2_better",
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemKlauspostS2(s2.WriterBetterCompression())
				},
			},
			{
				Name: "s2_best",
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemKlauspostS2(s2.WriterBestCompression())
				},
			},
		}...)
	}

	if os.Getenv(BENCH_S2_SNAPPY) != "" {
		compressionSystems = append(compressionSystems, []TestCompressionSystem{
			{
				Name: "s2_snappy_compat_default",
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemKlauspostS2(s2.WriterSnappyCompat())
				},
			},
			{
				Name: "s2_snappy_compat_better",
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemKlauspostS2(s2.WriterBetterCompression(), s2.WriterSnappyCompat())
				},
			},
			{
				Name: "s2_snappy_compat_best",
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemKlauspostS2(s2.WriterBestCompression(), s2.WriterSnappyCompat())
				},
			},
		}...)
	}

	if benchAll || os.Getenv(BENCH_DEFLATE) != "" {
		compressionSystems = append(compressionSystems, []TestCompressionSystem{
			{
				Name: "deflate_fastest",
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemKlauspostDeflate(ssb_compression.KlauspostDeflateLevelBestSpeed)
				},
			},
			{
				Name: "deflate_best",
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemKlauspostDeflate(ssb_compression.KlauspostDeflateLevelBestCompression)
				},
			},
			{
				Name: "deflate_default",
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemKlauspostDeflate(ssb_compression.KlauspostDeflateLevelDefaultCompression)
				},
			},
		}...)
	}

	if benchAll || os.Getenv(BENCH_LZMA) != "" {
		compressionSystems = append(compressionSystems, []TestCompressionSystem{
			{
				Name: "lzma",
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemLZMA()
				},
			},
			{
				Name: "lzma2",
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemLZMA2()
				},
			},
		}...)
	}

	if benchAll || os.Getenv(BENCH_ZSTD) != "" {
		for _, v := range []zstd.EncoderLevel{zstd.SpeedFastest, zstd.SpeedDefault, zstd.SpeedBestCompression} {
			level := v
			compressionSystems = append(compressionSystems, TestCompressionSystem{
				Name: fmt.Sprintf("zstd_%02d", level),
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemKlauspostZSTD(level)
				},
			})
		}
	}

	if benchAll || os.Getenv(BENCH_BROTLI) != "" {
		for _, v := range []int{brotli.DefaultCompression, brotli.BestCompression, brotli.BestSpeed} {
			level := v
			compressionSystems = append(compressionSystems, TestCompressionSystem{
				Name: fmt.Sprintf("brotli_%02d", level),
				Create: func() (ssb_compression.CompressionSystem, error) {
					return ssb_compression.NewSystemBrotli(level)
				},
			})
		}
	}

	for _, messageSource := range messageSources {
		b.Run(messageSource.Name, func(b *testing.B) {
			messages, err := Load(testdataFilepath(messageSource.Path))
			require.NoError(b, err)

			rand.Shuffle(len(messages), func(i, j int) {
				messages[i], messages[j] = messages[j], messages[i]
			})

			messages = messages[:int(float64(len(messages))*messageSource.SampleSize)]

			b.Log("source:", messageSource.Name, "number of messages:", len(messages))

			for _, batchN := range batches {
				b.Run(fmt.Sprintf("batch_%d", batchN), func(b *testing.B) {
					for _, system := range compressionSystems {
						s, err := system.Create()
						require.NoError(b, err)

						b.Run(system.Name, func(b *testing.B) {
							uncompressedSize := 0
							compressedSize := 0

							buf := &bytes.Buffer{}

							for i := 0; i < b.N; i++ {
								uncompressedSize = 0
								compressedSize = 0

								j := 0
								for {
									buf.Reset()

									for k := j; k < j+batchN; k++ {
										if k >= len(messages)-1 {
											break
										}
										uncompressedSize += len(messages[k].Payload)
										buf.Write(messages[k].Payload)
									}

									out, err := s.Compress(buf.Bytes())
									if err != nil {
										b.Fatal(err)
									}

									compressedSize += len(out)

									j += batchN
									if j >= len(messages)-1 {
										break
									}
								}
							}

							b.ReportMetric(float64(uncompressedSize), "bytes_per_messages")
							b.ReportMetric(float64(compressedSize), "compressedbytes_per_messages")
							b.ReportMetric(float64(uncompressedSize)/float64(compressedSize), "ratio")
						})
					}
				})
			}
		})
	}
}

type MessageType int

const (
	MessageTypeSent MessageType = iota
	MessageTypeReceived
)

type Message struct {
	Payload     []byte
	MessageType MessageType
}

func Load(filepath string) ([]Message, error) {
	var messages []Message

	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(nil, 50*bufio.MaxScanTokenSize)

	for scanner.Scan() {
		line := scanner.Text()

		msgType, err := messageType(line)
		if err != nil {
			panic(err)
		}

		line = strings.TrimPrefix(line, prefixSend)
		line = strings.TrimPrefix(line, prefixReceive)
		line = strings.TrimSpace(line)

		payload, err := hex.DecodeString(line)
		if err != nil {
			panic(err)
		}

		messages = append(messages, Message{
			Payload:     payload,
			MessageType: msgType,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return messages, nil
}

const (
	prefixSend    = "hexsend"
	prefixReceive = "hexreceive"
)

func messageType(line string) (MessageType, error) {
	if strings.HasPrefix(line, prefixSend) {
		return MessageTypeSent, nil
	}

	if strings.HasPrefix(line, prefixReceive) {
		return MessageTypeReceived, nil
	}

	return 0, errors.New("unknown")
}

func testdataFilepath(s string) string {
	return path.Join(testdataDirectory(), s)
}

func testdataDirectory() string {
	return path.Join(currentDirectory(), "testdata")
}

func currentDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}
