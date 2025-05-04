package pipeline

import (
	"log"
	"os"

	"github.com/babaYaga451/go-tnt-automation/internal/apiclient"
	"github.com/babaYaga451/go-tnt-automation/internal/discover"
	"github.com/babaYaga451/go-tnt-automation/internal/enricher"
	"github.com/babaYaga451/go-tnt-automation/internal/model"
	"github.com/babaYaga451/go-tnt-automation/internal/sampler"
)

func RunPipeLine(inputDir, mapFile, apiURL string, k, workers int) <-chan model.TestResult {

	log.Println("Starting pipeline")
	log.Println("Using inputDir:", inputDir)
	log.Println("Using mapFile:", mapFile)
	destMap, err := model.LoadDestInfo(mapFile)
	if err != nil {
		os.Exit(1)
	}

	fileCh := discover.DiscoverFiles(inputDir)
	sampleCh := sampler.SampleStage(fileCh, k, workers)
	enrichCh := enricher.EnrichStage(sampleCh, destMap, workers)
	resultCh := apiclient.APIStage(enrichCh, apiURL, workers)

	return resultCh
}
