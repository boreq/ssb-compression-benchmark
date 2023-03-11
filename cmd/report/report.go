package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/boreq/errors"
	"github.com/boreq/ssb-compression-benchmark/report"
	gochart "github.com/wcharczuk/go-chart/v2"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	results, err := report.GetBenchResults(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "error getting bench results")
	}

	directory := path.Join(
		"results",
		fmt.Sprintf("%s-%s-%s", results.Cpu, results.Goarch, results.Goos),
	)

	if err := os.RemoveAll(directory); err != nil {
		return errors.Wrap(err, "error removing directory")
	}

	if err := os.MkdirAll(directory, 0700); err != nil {
		return errors.Wrap(err, "error recreating directory")
	}

	readmeBuffer := bytes.NewBuffer(nil)
	readmeBuffer.WriteString("# Results\n")
	readmeBuffer.WriteString("```\n")
	readmeBuffer.WriteString(fmt.Sprintf("goarch=%s\n", results.Goarch))
	readmeBuffer.WriteString(fmt.Sprintf("goos=%s\n", results.Goos))
	readmeBuffer.WriteString(fmt.Sprintf("cpu=%s\n", results.Cpu))
	readmeBuffer.WriteString("```\n")

	readmeBuffer.WriteString("## Performance\n")

	for _, result := range results.CompressionResults {
		resultsChart, err := report.MakeResultChart(result)
		if err != nil {
			return errors.Wrap(err, "error creating chart")
		}

		filename := fmt.Sprintf(
			"%s.png",
			strings.Replace(result.BenchmarkName, string(os.PathSeparator), "-", -1),
		)

		f, err := os.Create(path.Join(directory, filename))
		if err != nil {
			return errors.Wrap(err, "error creating chart file")
		}

		if err := resultsChart.Render(gochart.PNG, f); err != nil {
			return errors.Wrap(err, "error rendering the chart")
		}

		readmeBuffer.WriteString(fmt.Sprintf("### %s\n", result.BenchmarkName))
		readmeBuffer.WriteString(fmt.Sprintf("![](./%s)\n", mdSafe(filename)))
		readmeBuffer.WriteString("```\n")
		sort.Slice(result.Systems, func(i, j int) bool {
			return result.Systems[i].Ratio > result.Systems[j].Ratio
		})
		for _, system := range result.Systems {
			readmeBuffer.WriteString(fmt.Sprintf("%21s = %.2f ratio\n", system.SystemName, system.Ratio))
		}
		readmeBuffer.WriteString("```\n")

	}

	readmeFile, err := os.Create(path.Join(directory, "README.md"))
	if err != nil {
		return errors.Wrap(err, "error creating readme")
	}

	if _, err := readmeBuffer.WriteTo(readmeFile); err != nil {
		return errors.Wrap(err, "error writing to readme file")
	}

	return nil
}

func mdSafe(s string) string {
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	s = strings.ReplaceAll(s, " ", "&#32;")
	return s
}
