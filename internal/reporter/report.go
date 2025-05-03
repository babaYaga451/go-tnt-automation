package report

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/babaYaga451/go-tnt-automation/internal/model"
)

// CollectResults gathers all TestResults from the channel
func CollectResults(resultCh <-chan model.TestResult) []model.TestResult {
	var all []model.TestResult
	for r := range resultCh {
		all = append(all, r)
	}
	return all
}

// WriteJUnit writes a JUnit XML report based on the results
func WriteJUnit(path string, results []model.TestResult) error {
	startTime := time.Now()
	grouped := make(map[string][]model.TestResult)
	for _, r := range results {
		key := fmt.Sprintf("%s.%s", r.Record.Shipper, r.Record.Origin)
		grouped[key] = append(grouped[key], r)
	}

	var allSuites []model.Testsuite
	totalTests := 0
	totalFailures := 0

	for key, res := range grouped {
		suite := model.Testsuite{
			Name:      key,
			File:      "transit-api.go",
			Timestamp: time.Now().Format(time.RFC3339),
		}

		for _, r := range res {
			tc := model.Testcase{
				Classname:  key,
				Name:       fmt.Sprintf("%s -> %s", r.Record.Origin, r.Record.Destination),
				Time:       fmt.Sprintf("%.3f", r.Duration.Seconds()),
				Assertions: 1,
			}

			if r.Err != nil {
				tc.Failure = &model.Failure{
					Message: r.Err.Error(),
					Type:    "APIError",
					Content: fmt.Sprintf("ActualDays: %d, Expected: %d", r.ActualDays, r.Record.TransitDays),
				}
				suite.Failures++
			} else if r.ActualDays != r.Record.TransitDays {
				tc.Failure = &model.Failure{
					Message: "Transit mismatch",
					Content: fmt.Sprintf("Expected %d, got %d", r.Record.TransitDays, r.ActualDays),
				}
				suite.Failures++
			}

			suite.Testcases = append(suite.Testcases, tc)
		}

		suite.Tests = len(suite.Testcases)
		allSuites = append(allSuites, suite)
		totalTests += suite.Tests
		totalFailures += suite.Failures
	}

	final := model.Testsuites{
		Name:      "Transit Tests",
		Tests:     totalTests,
		Failures:  totalFailures,
		Timestamp: time.Now().Format(time.RFC3339),
		Time:      fmt.Sprintf("%.3f", float64(time.Since(startTime).Milliseconds())/1000),
		Suites:    allSuites,
	}

	// Write XML
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := xml.NewEncoder(file)
	enc.Indent("", "  ")
	return enc.Encode(final)
}

func CountFailures(results []model.TestResult) int {
	failures := 0
	for _, r := range results {
		if r.Err != nil || r.ActualDays != r.Record.TransitDays {
			failures++
		}
	}
	return failures
}
