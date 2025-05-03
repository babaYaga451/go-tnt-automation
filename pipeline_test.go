package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/babaYaga451/go-tnt-automation/internal/pipeline"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

func TestTransitResults(t *testing.T) {
	results := pipeline.RunPipeLine(
		"./data",
		"dest.csv",
		"http://localhost:8080/transit",
		10,
		4,
	)

	for _, tr := range results {
		testName := fmt.Sprintf("%s → %s", tr.Record.Origin, tr.Record.Destination)

		runner.Run(t, testName, func(t provider.T) {
			t.Labels(
				allure.ParentSuiteLabel("TestTransitResults"),
				allure.SuiteLabel(tr.Record.Shipper),
				allure.SubSuiteLabel(fmt.Sprintf("transitDays:%d", tr.Record.TransitDays)),

				allure.TagLabel("shipper:"+tr.Record.Shipper),
				allure.TagLabel("origin:"+tr.Record.Origin),
				allure.TagLabel("destination:"+tr.Record.Destination),
				allure.TagLabel("zip:"+tr.Record.Zip),
				allure.TagLabel("city:"+strings.ToLower(tr.Record.City)),
				allure.TagLabel("state:"+strings.ToLower(tr.Record.State)),
				allure.TagLabel(fmt.Sprintf("transitDays:%d", tr.Record.TransitDays)),
			)

			t.WithNewStep("Transit Check", func(s provider.StepCtx) {
				s.WithAttachments(
					allure.NewAttachment("Origin", allure.Text, []byte(tr.Record.Origin)),
					allure.NewAttachment("Destination", allure.Text, []byte(tr.Record.Destination)),
				)

				if tr.Err != nil {
					s.NewStep("API call failed: " + tr.Err.Error())
					s.Fail()
					t.Errorf("API call failed: %v", tr.Err)
				} else if tr.ActualDays != tr.Record.TransitDays {
					msg := fmt.Sprintf("Transit mismatch: expected %d, got %d", tr.Record.TransitDays, tr.ActualDays)
					s.NewStep(msg)
					s.Fail()
					t.Error(msg)
				}
			})
		})
	}
}
