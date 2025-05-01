package sampler

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/babaYaga451/go-tnt-automation/internal/model"
)

// reservoir holds up to k items and randomly replaces them
// as more items stream through.
type reservoir struct {
	k     int
	items []model.Record
	count int
}

// newReservoir creates a new reservoir that will keep k samples.
func newReservoir(k int) *reservoir {
	return &reservoir{k: k, items: make([]model.Record, 0, k)}
}

// add considers rec for inclusion in the reservoir.
// If we haven’t filled k slots yet, we append.
// Once full, we pick a random index in [0, count)
// and if that index is < k, we replace items[index].
func (r *reservoir) add(rec model.Record) {
	r.count++
	if len(r.items) < r.k {
		r.items = append(r.items, rec)
	} else {
		// pick a random int in [0, r.count)
		idx := rand.Intn(r.count)
		if idx < r.k {
			r.items[idx] = rec
		}
	}
}

// SampleStage does reservoir sampling (k per transit-day) across files
func SampleStage(fileCh <-chan string, k, workers int) <-chan model.Record {
	out := make(chan model.Record, k*workers)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			samplers := map[int]*reservoir{}
			for path := range fileCh {

				f, err := os.Open(path)
				if err != nil {
					continue
				}
				fileName := filepath.Base(path)
				shipper := strings.TrimSuffix(fileName, filepath.Ext(fileName))

				rdr := csv.NewReader(f)
				rdr.Comma = '|'
				rdr.FieldsPerRecord = -1
				rdr.Read() // skip header

				for {
					rec, err := rdr.Read()
					if err == io.EOF {
						break
					}
					if err != nil || len(rec) < 4 {
						continue
					}
					days, e := strconv.Atoi(rec[3])
					if e != nil {
						continue
					}
					r := model.Record{
						Origin:      rec[0],
						Destination: rec[1],
						TransitDays: days,
						Shipper:     shipper,
					}
					if _, ok := samplers[days]; !ok {
						samplers[days] = newReservoir(k)
					}
					fmt.Println(r)
					samplers[days].add(r)
				}
				f.Close()
				for _, rs := range samplers {
					for _, item := range rs.items {
						out <- item
					}
				}
				samplers = map[int]*reservoir{}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
