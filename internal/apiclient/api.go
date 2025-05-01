package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/babaYaga451/go-tnt-automation/internal/model"
)

// APIStage executes calls and validates transit days
func APIStage(enrichCh <-chan model.Record, apiURL string, workers int) <-chan model.TestResult {
	out := make(chan model.TestResult, cap(enrichCh))
	client := &http.Client{Timeout: 10 * time.Second}
	type apiResponse struct {
		TransitDays int `json:"transitDays"`
	}
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rec := range enrichCh {
				start := time.Now()
				url := fmt.Sprintf("%s?origin=%s&dest=%s&shipper=%s&city=%s&state=%s&zip=%s",
					apiURL, rec.Origin, rec.Destination, rec.Shipper, rec.City, rec.State, rec.Zip)
				tr := model.TestResult{Record: rec}
				resp, err := client.Get(url)
				if err != nil {
					tr.Err = err
				} else {
					var r apiResponse
					err := json.NewDecoder(resp.Body).Decode(&r)
					resp.Body.Close()
					if err != nil {
						tr.Err = err
					} else {
						tr.ActualDays = r.TransitDays
					}
				}
				tr.Duration = time.Since(start)
				out <- tr
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
