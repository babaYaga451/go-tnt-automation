package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/babaYaga451/go-tnt-automation/internal/pipeline"
	report "github.com/babaYaga451/go-tnt-automation/internal/reporter"
)

var (
	inputDir  = flag.String("inputDir", "./data", "Directory containing .txt files (each named <shipper>.txt)")
	mapFile   = flag.String("mapFile", "dest.csv", "CSV mapping destination→city,state,zip")
	apiURL    = flag.String("apiURL", "http://localhost:8080/transit", "API endpoint URL")
	k         = flag.Int("k", 10, "Samples per transit-day group per file")
	workers   = flag.Int("workers", runtime.NumCPU(), "Parallel workers per stage")
	outputDir = flag.String("outputFile", "results.xml", "Junit XML path")
)

func main() {
	flag.Parse()
	results := pipeline.RunPipeLine(*inputDir, *mapFile, *apiURL, *k, *workers)
	if err := report.WriteJUnit(*outputDir, results); err != nil {
		os.Exit(1)
	}

	fmt.Printf("Done: %d tests, %d failures\n", len(results), report.CountFailures(results))
}
