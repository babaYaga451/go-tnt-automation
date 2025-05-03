package model

import (
	"encoding/csv"
	"encoding/xml"
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

type Testsuites struct {
	XMLName   xml.Name    `xml:"testsuites"`
	Name      string      `xml:"name,attr,omitempty"`
	Tests     int         `xml:"tests,attr"`
	Failures  int         `xml:"failures,attr,omitempty"`
	Errors    int         `xml:"errors,attr,omitempty"`
	Skipped   int         `xml:"skipped,attr,omitempty"`
	Time      string      `xml:"time,attr,omitempty"`
	Timestamp string      `xml:"timestamp,attr,omitempty"`
	Suites    []Testsuite `xml:"testsuite"`
}

type Testsuite struct {
	Name       string     `xml:"name,attr"`
	Tests      int        `xml:"tests,attr"`
	Failures   int        `xml:"failures,attr"`
	Errors     int        `xml:"errors,attr"`
	Skipped    int        `xml:"skipped,attr"`
	Assertions int        `xml:"assertions,attr,omitempty"`
	Time       string     `xml:"time,attr,omitempty"`
	Timestamp  string     `xml:"timestamp,attr,omitempty"`
	File       string     `xml:"file,attr,omitempty"`
	Properties []Property `xml:"properties>property,omitempty"`
	SystemOut  string     `xml:"system-out,omitempty"`
	SystemErr  string     `xml:"system-err,omitempty"`
	Testcases  []Testcase `xml:"testcase"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Testcase struct {
	Name       string     `xml:"name,attr"`
	Classname  string     `xml:"classname,attr"`
	Assertions int        `xml:"assertions,attr,omitempty"`
	Time       string     `xml:"time,attr,omitempty"`
	File       string     `xml:"file,attr,omitempty"`
	Line       int        `xml:"line,attr,omitempty"`
	Failure    *Failure   `xml:"failure,omitempty"`
	Error      *Error     `xml:"error,omitempty"`
	Skipped    *Skipped   `xml:"skipped,omitempty"`
	SystemOut  string     `xml:"system-out,omitempty"`
	SystemErr  string     `xml:"system-err,omitempty"`
	Properties []Property `xml:"properties>property,omitempty"`
}

type Failure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr,omitempty"`
	Content string `xml:",chardata"`
}

type Error struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr,omitempty"`
	Content string `xml:",chardata"`
}

type Skipped struct {
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
