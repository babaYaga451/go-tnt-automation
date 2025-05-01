package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Response is the JSON structure returned by the API.
type Response struct {
	TransitDays int `json:"transitDays"`
}

func transitHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	fmt.Println(r.URL.String())
	originStr := r.URL.Query().Get("origin")
	destStr := r.URL.Query().Get("dest")

	// Compute a deterministic transitDays (1–10) based on numeric origin/dest
	origin, err1 := strconv.Atoi(originStr)
	dest, err2 := strconv.Atoi(destStr)
	var days int
	if err1 == nil && err2 == nil {
		days = int(math.Abs(float64(dest-origin)))%10 + 1
	} else {
		days = rand.Intn(10) + 1
	}

	// Simulate network latency: 100–500 ms
	delay := time.Duration(rand.Intn(400)+100) * time.Millisecond
	time.Sleep(delay)

	// Build and send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{TransitDays: days})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Route setup
	http.HandleFunc("/transit", transitHandler)

	// Start server on port 8080
	log.Println("Transit API server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
