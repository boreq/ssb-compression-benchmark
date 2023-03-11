package report

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/boreq/errors"
	"github.com/wcharczuk/go-chart/v2"
)

type BenchResults struct {
	Goos               string
	Goarch             string
	Cpu                string
	CompressionResults []CompressionBenchResult
}

type CompressionBenchResult struct {
	BenchmarkName string
	Systems       []SystemPerformanceBenchResult
}

type SystemPerformanceBenchResult struct {
	SystemName string
	Ratio      float64
}

func GetBenchResults(r io.Reader) (BenchResults, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return BenchResults{}, errors.Wrap(err, "error reading all")
	}

	var result BenchResults

	scan := bufio.NewScanner(bytes.NewReader(b))
	for scan.Scan() {
		if err := parseLine(scan.Text(), &result); err != nil {
			return BenchResults{}, errors.Wrap(err, "error parsing line")
		}
	}

	if err := scan.Err(); err != nil {
		return BenchResults{}, errors.Wrap(err, "scan error")
	}

	if result.Cpu == "" || result.Goarch == "" || result.Goos == "" {
		return BenchResults{}, fmt.Errorf("missing execution environment info in output: '%+v'", result)
	}

	results, err := getBenchResults(bytes.NewReader(b))
	if err != nil {
		return BenchResults{}, errors.Wrap(err, "error getting performance results")
	}

	result.CompressionResults = results

	return result, err
}

const lineSep = ":"

func parseLine(line string, result *BenchResults) error {
	splitLine := strings.SplitN(line, lineSep, 2)
	if len(splitLine) != 2 {
		return nil
	}

	key := splitLine[0]
	value := strings.TrimSpace(splitLine[1])

	switch key {
	case "goos":
		result.Goos = value
	case "goarch":
		result.Goarch = value
	case "cpu":
		result.Cpu = value
	case "pkg":
	default:
		return errors.New("unknown line")
	}
	return nil
}

func getBenchResults(r io.Reader) ([]CompressionBenchResult, error) {
	var results []CompressionBenchResult

	set, err := ParseSet(r)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing set")
	}

	for _, benchmarks := range set {
		for _, benchmark := range benchmarks {
			systemName, benchmarkName, err := ParseBenchmarkName(benchmark.Name)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing benchmark name")
			}

			bench, ok := findBenchmark(results, benchmarkName)
			if !ok {
				results = append(results, CompressionBenchResult{
					BenchmarkName: benchmarkName,
					Systems:       nil,
				})
				bench = &results[len(results)-1]
			}

			bench.Systems = append(bench.Systems, SystemPerformanceBenchResult{
				SystemName: systemName,
				Ratio:      benchmark.Measurements["ratio"],
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].BenchmarkName < results[j].BenchmarkName
	})

	for _, result := range results {
		sort.Slice(result.Systems, func(i, j int) bool {
			return result.Systems[i].SystemName < result.Systems[j].SystemName
		})
	}

	return results, nil
}

const (
	chartHeight   = 800
	chartWidth    = 2000
	chartBarWidth = 100
)

func MakeResultChart(result CompressionBenchResult) (chart.BarChart, error) {
	graph := chart.BarChart{
		Title: result.BenchmarkName,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:       chartHeight,
		BarWidth:     chartBarWidth,
		Width:        chartWidth,
		UseBaseValue: true,
		BaseValue:    1.0,
		YAxis: chart.YAxis{
			Name:      "ratio",
			NameStyle: chart.StyleTextDefaults(),
			Range: &chart.ContinuousRange{
				Min: 0.9,
				Max: 2.5,
			},
			Ticks: []chart.Tick{
				{
					Value: 1,
					Label: "1",
				},
				{
					Value: 1.5,
					Label: "1.5",
				},
				{
					Value: 2,
					Label: "2",
				},
				{
					Value: 2.5,
					Label: "2.5",
				},
			},
		},
		Elements: []chart.Renderable{
			func(r chart.Renderer, canvasBox chart.Box, defaults chart.Style) {
				defaults.FontColor = drawing.ColorBlack
				defaults.GetTextOptions().WriteToRenderer(r)
				r.SetTextRotation(3.14 / 2)
				text := "Compression ratio"
				r.Text(text, canvasBox.Width()+60, canvasBox.Height()/2)
			},
		},
	}

	var max *float64
	for _, system := range result.Systems {
		if max == nil || *max < system.Ratio {
			tmp := system.Ratio
			max = &tmp
		}
	}

	var min *float64
	for _, system := range result.Systems {
		if min == nil || *min > system.Ratio {
			tmp := system.Ratio
			min = &tmp
		}
	}

	for _, system := range result.Systems {
		value := chart.Value{
			Label: system.SystemName,
			Value: system.Ratio,
		}

		baselineTransparency := 10.0

		if system.Ratio < 1 {
			fractionFromOneToMin := (1 - system.Ratio) / (1 - *min)
			color := drawing.Color{R: 231, G: 76, B: 60, A: uint8(baselineTransparency + (255.0-baselineTransparency)*fractionFromOneToMin)}

			redStyle := chart.Style{
				FillColor:   color,
				StrokeColor: color,
				StrokeWidth: 0.01,
			}
			value.Style = redStyle
		} else {
			fractionFromOneToMax := (system.Ratio - 1) / (*max - 1)
			color := drawing.Color{R: 46, G: 204, B: 113, A: uint8(baselineTransparency + (255.0-baselineTransparency)*fractionFromOneToMax)}

			gradientStyle := chart.Style{
				FillColor:   color,
				StrokeColor: color,
				StrokeWidth: 0.01,
			}
			value.Style = gradientStyle
		}

		graph.Bars = append(graph.Bars, value)
	}

	return graph, nil
}

var systemNames = map[string]string{
	"brotli_06-":                "Brotli (best)",
	"brotli_11-":                "Brotli (default)",
	"brotli_00-":                "Brotli (fastest)",
	"deflate_best-":             "Deflate (best)",
	"deflate_fastest-":          "Deflate (fastest)",
	"deflate_default-":          "Deflate (default)",
	"lzma-":                     "LZMA",
	"lzma2-":                    "LZMA2",
	"s2_best-":                  "S2 (best)",
	"s2_better-":                "S2 (better)",
	"s2_default-":               "S2 (default)",
	"s2_snappy_compat_best-":    "S2 with Snappy compatibility (best)",
	"s2_snappy_compat_default-": "S2 with Snappy compatibility (default)",
	"s2_snappy_compat_better-":  "S2 with Snappy compatibility (better)",
	"snappy-":                   "Snappy",
	"zstd_01-":                  "ZSTD (fastest)",
	"zstd_02-":                  "ZSTD (default)",
	"zstd_04-":                  "ZSTD (best)",
}

var benchmarkNames = map[string]string{
	"BenchmarkLines.many_feed_messages.batch_1":   "Many feed messages (compressing messages individually)",
	"BenchmarkLines.many_feed_messages.batch_10":  "Many feed messages (compressing messages in batches of 10)",
	"BenchmarkLines.many_feed_messages.batch_100": "Many feed messages (compressing messages in batches of 100)",

	"BenchmarkLines.few_feed_messages.batch_1":   "Few feed messages (compressing messages individually)",
	"BenchmarkLines.few_feed_messages.batch_10":  "Few feed messages (compressing messages in batches of 10)",
	"BenchmarkLines.few_feed_messages.batch_100": "Few feed messages (compressing messages in batches of 100)",
}

func ParseBenchmarkName(name string) (string, string, error) {
	split := strings.Split(name, "/")
	if len(split) != 4 {
		return "", "", errors.New("invalid name")
	}

	found := false
	systemName := split[3]
	for prefix, replacement := range systemNames {
		if strings.HasPrefix(systemName, prefix) {
			systemName = replacement
			found = true
			break
		}
	}
	if !found {
		return "", "", fmt.Errorf("unknown system name '%s'", systemName)
	}

	found = false
	benchmarkName := strings.Join(split[:3], ".")
	for name, replacement := range benchmarkNames {
		if benchmarkName == name {
			benchmarkName = replacement
			found = true
			break
		}
	}
	if !found {
		return "", "", fmt.Errorf("unknown benchmark name '%s'", benchmarkName)
	}

	return systemName, benchmarkName, nil
}

func findBenchmark(results []CompressionBenchResult, benchmarkName string) (*CompressionBenchResult, bool) {
	for i := range results {
		if results[i].BenchmarkName == benchmarkName {
			return &results[i], true
		}
	}
	return nil, false
}

// Benchmark is one run of a single benchmark.
type Benchmark struct {
	Name         string             // benchmark name
	N            int                // number of iterations
	Measurements map[string]float64 // unit => value
	Ord          int                // ordinal position within a benchmark run
}

// ParseLine extracts a Benchmark from a single line of testing.B
// output.
func ParseLine(line string) (*Benchmark, error) {
	fields := strings.Fields(line)

	// Two required, positional fields: Name and iterations.
	if len(fields) < 2 {
		return nil, fmt.Errorf("two fields required, have %d", len(fields))
	}
	if !strings.HasPrefix(fields[0], "Benchmark") {
		return nil, fmt.Errorf(`first field does not start with "Benchmark"`)
	}
	n, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, err
	}
	b := &Benchmark{Name: fields[0], N: n, Measurements: make(map[string]float64)}

	// Parse any remaining pairs of fields; we've parsed one pair already.
	for i := 1; i < len(fields)/2; i++ {
		if err := b.parseMeasurement(fields[i*2], fields[i*2+1]); err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (b *Benchmark) parseMeasurement(quant string, unit string) error {
	f, err := strconv.ParseFloat(quant, 64)
	if err != nil {
		return fmt.Errorf("error parsing quantity '%s': %w", quant, err)
	}
	b.Measurements[unit] = f
	return nil
}

func (b *Benchmark) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s %d", b.Name, b.N)
	for unit, quant := range b.Measurements {
		fmt.Fprintf(buf, " %.2f %s", quant, unit)
	}
	return buf.String()
}

// Set is a collection of benchmarks from one
// testing.B run, keyed by name to facilitate comparison.
type Set map[string][]*Benchmark

// ParseSet extracts a Set from testing.B output.
// ParseSet preserves the order of benchmarks that have identical
// names.
func ParseSet(r io.Reader) (Set, error) {
	bb := make(Set)
	scan := bufio.NewScanner(r)
	ord := 0
	for scan.Scan() {
		if b, err := ParseLine(scan.Text()); err == nil {
			b.Ord = ord
			ord++
			bb[b.Name] = append(bb[b.Name], b)
		}
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}

	return bb, nil
}
