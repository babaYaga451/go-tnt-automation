package enricher

import (
	"sync"

	"github.com/babaYaga451/go-tnt-automation/internal/model"
)

// EnrichStage looks up city/state/zip and tags shipper
func EnrichStage(sampleCh <-chan model.Record, zipMap map[string]model.DestInfo, workers int) <-chan model.Record {

	out := make(chan model.Record, cap(sampleCh))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rec := range sampleCh {
				if info, ok := zipMap[rec.Destination]; ok {
					rec.City = info.City
					rec.State = info.State
					rec.Zip = info.Zip
				}
				out <- rec
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
