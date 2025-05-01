package model

import (
	"encoding/csv"
	"io"
	"os"
	"time"
)

type Record struct {
	Origin      string
	Destination string
	City        string
	State       string
	Zip         string
	Shipper     string
	TransitDays int
}

type TestResult struct {
	Record     Record
	ActualDays int
	Err        error
	Duration   time.Duration
}

type DestInfo struct {
	City  string
	State string
	Zip   string
}

type Testsuite struct {
	Name     string     `xml:"name,attr"`
	Tests    int        `xml:"tests,attr"`
	Failures int        `xml:"failures,attr"`
	Cases    []Testcase `xml:"testcase"`
}

type Testcase struct {
	Class   string   `xml:"classname,attr"`
	Name    string   `xml:"name,attr"`
	Time    string   `xml:"time,attr"`
	Failure *Failure `xml:"failure,omitempty"`
}

type Failure struct {
	Message string `xml:"message,attr"`
}

func LoadDestInfo(path string) (map[string]DestInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rdr := csv.NewReader(f)
	// If you have a header row, call rdr.Read() here to skip it:
	// _, _ = rdr.Read()

	destMap := make(map[string]DestInfo)
	for {
		rec, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(rec) < 4 {
			// skip bad lines
			continue
		}
		code := rec[0]
		destMap[code] = DestInfo{
			City:  rec[1],
			State: rec[2],
			Zip:   rec[3],
		}
	}
	return destMap, nil
}
