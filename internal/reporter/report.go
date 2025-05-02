package report

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

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
	suite := model.Testsuite{
		Name:  "TransitAPITests",
		Tests: len(results),
	}

	failCount := 0

	for _, r := range results {
		tc := model.Testcase{
			Class: fmt.Sprintf("%s.%s", r.Record.Shipper, r.Record.Origin),
			Name:  fmt.Sprintf("%s -> %s", r.Record.Origin, r.Record.Destination),
			Time:  fmt.Sprintf("%.3f", r.Duration.Seconds()), // JUnit expects seconds as float
		}

		if r.Err != nil || r.ActualDays != r.Record.TransitDays {
			failCount++
			var msg string
			if r.Err != nil {
				msg = r.Err.Error()
			} else {
				msg = fmt.Sprintf(
					"Mismatch: Shipper=%s, Origin=%s, Destination=%s, Expected=%d, Actual=%d",
					r.Record.Shipper,
					r.Record.Origin,
					r.Record.Destination,
					r.Record.TransitDays,
					r.ActualDays,
				)
			}
			tc.Failure = &model.Failure{
				Message: msg,
			}
		}

		suite.Cases = append(suite.Cases, tc)
	}

	suite.Failures = failCount

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")
	return enc.Encode(suite)
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
