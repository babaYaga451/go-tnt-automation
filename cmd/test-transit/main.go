package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/babaYaga451/go-tnt-automation/internal/apiclient"
	"github.com/babaYaga451/go-tnt-automation/internal/discover"
	"github.com/babaYaga451/go-tnt-automation/internal/enricher"
	"github.com/babaYaga451/go-tnt-automation/internal/model"
	report "github.com/babaYaga451/go-tnt-automation/internal/reporter"
	"github.com/babaYaga451/go-tnt-automation/internal/sampler"
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
	// flags omitted for brevity
	flag.Parse() // ✅ This is critical

	log.Println("Starting pipeline")
	log.Println("Using inputDir:", *inputDir)
	log.Println("Using mapFile:", *mapFile)
	destMap, err := model.LoadDestInfo(*mapFile)
	if err != nil {
		os.Exit(1)
	}

	fileCh := discover.DiscoverFiles(*inputDir)
	sampleCh := sampler.SampleStage(fileCh, *k, *workers)
	enrichCh := enricher.EnrichStage(sampleCh, destMap, *workers)
	resultCh := apiclient.APIStage(enrichCh, *apiURL, *workers)

	results := report.CollectResults(resultCh)
	if err := report.WriteJUnit(*outputDir, results); err != nil {
		os.Exit(1)
	}
	fmt.Printf("Done: %d tests, %d failures\n", len(results), report.CountFailures(results))
}
