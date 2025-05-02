# go-tnt-automation

Automation to validate transit days by calling api service

# 🚀 go-tnt-tester

A high-performance test automation tool written in Go to validate transit times for origin–destination ZIP pairs at scale.

This CLI utility:

- Processes 1000s of `.txt` files with transit expectations
- Samples test cases by transit days
- Calls an external `/transit` API
- Validates results against expected values
- Outputs structured `JUnit` reports (compatible with Jenkins + Allure)
- Supports parallel sharding in CI pipelines

---

## 🧱 Features

- ⚡ Multi-threaded and memory efficient (streaming, batching)
- 📂 Automatically sharded file input
- 🧪 Outputs JUnit XML reports per shard
- 🐳 Docker-compatible
- ⚙️ Easy CI/CD integration (Jenkins, GitHub Actions, etc.)
- 📊 Supports Allure reports (via converter)

---

## 📁 Input File Format

```text
origin|destination|state|transitDays
10001|90210|CA|3
10001|10010|NY|2
...
```

# 🧠 Understanding the go-tnt-automation Pipeline (Assembly Line Metaphor)

Think of your program like a team of helpers on an assembly line. Each helper has one simple job and works in its own space (folder), so they don’t get in each other’s way.

Here’s how it works, step by step:

---

## 1. `cmd/transit-tester/main.go` – 🧑‍💼 The Boss

- This is the **"boss"** who kicks everything off:

  - Collects all user settings: input directory, API URL, shipper, etc.
  - Starts each helper (pipeline stage) and wires them together using channels.

---

## 2. `internal/discover/discover.go` – 🔍 The Scout

- The scout’s job is to **find all `.txt` files** in your folder.
- It walks the directory, and for each `*.txt` file it finds, it drops the filepath into a **channel** for the next stage to use.

---

## 3. `internal/sample/sample.go` – 🍬 The Sampler

- Imagine a huge pile of candy — but you only want a few of each color.
- The sampler reads each text file **line by line**, groups by `transitDay`, and **keeps only `k` random lines** per group using reservoir sampling.
- Selected lines are sent to the next stage via a channel.

---

## 4. `internal/enrich/enrich.go` – 🎨 The Decorator

- Takes the sampled records (just origin, destination, transitDay) and looks up:

  - city
  - state
  - zip code

- It also tags each record with the shipper name.
- These enriched records are passed along.

---

## 5. `internal/api/api.go` – 📞 The Caller & Checker

- For every enriched record, this helper calls your API:

  - "Hey API, how many days to deliver from X to Y with Z?"

- It compares the **actual API response** with the **expected transit day** from the file.
- Packages the result: success/failure, how long it took, any errors.

---

## 6. `internal/report/report.go` – 🧾 The Reporter

- Gathers all the results from the previous step.
- Writes a **JUnit-style XML report** that Jenkins/Allure/etc. can consume.
- Tallies pass/fail counts.

---

## 📦 How These Helpers Communicate

- Each helper uses a **Go channel** as its inbox and outbox.
- No helper holds everything in memory — they **stream** data as it's processed.
- If one stage gets slow (e.g. the API server), its output channel fills up, which naturally causes earlier stages to **pause**.

---

## 🔗 Example of Channels in Code

```go
fileCh   := make(chan string, 1000)
sampleCh := make(chan model.Record, 10000)
enrichCh := make(chan model.Record, 10000)
resultCh := make(chan model.TestResult, 64000)

// Discover stage writes to fileCh
// Sample stage reads from fileCh, writes to sampleCh
// Enrich stage reads from sampleCh, writes to enrichCh
// API stage reads from enrichCh, writes to resultCh
// Report stage reads from resultCh
```

Each `go func()` reads from its input channel and writes to its output. If the output channel is full, it blocks — causing natural back pressure.

---

## 📌 What is Back Pressure?

> Back pressure is a natural way for a pipeline to **slow down producers when consumers are overwhelmed**, preventing memory overflow and keeping throughput stable.

In Go, back pressure is implemented using **bounded channels**. If a consumer is slower than the producer, the channel fills up, causing the producer to block.

---

## 🧱 Pipeline Structure

The application consists of the following pipeline stages:

```
[discover] → [sample] → [enrich] → [api call] → [result writer]
```

Each stage communicates through a buffered channel:

```go
fileCh   := make(chan string, 1000)
sampleCh := make(chan Record, 10000)
enrichCh := make(chan Record, 10000)
resultCh := make(chan TestResult, 64000)
```

---

## 🔁 How Back Pressure Works

Let's say the API server is slow, and the `resultCh` becomes full:

1. The `api` stage tries to send to `resultCh` but **blocks** because it's full.
2. The API worker is now **paused**, and cannot receive from `enrichCh`.
3. This causes `enrichCh` to back up, blocking the `enrich` stage.
4. That blocks the `sample` stage, and so on — all the way back to `discover`.

### 🧩 Go makes this automatic:

- No need for semaphores or throttling logic
- Workers simply block on full channels

---

## 🧪 Code Snippet

```go
for rec := range enrichCh {
    ...
    resultCh <- tr // Blocks if resultCh is full → back pressure starts here
}
```

---

## ✅ Benefits

| Mechanism             | Benefit                                 |
| --------------------- | --------------------------------------- |
| Buffered channels     | Memory-safe concurrency                 |
| Goroutine blocking    | Naturally slows the producer            |
| Stage-by-stage pull   | Prevents overloading downstream systems |
| No polling/throttling | Pure Go concurrency model               |

---
