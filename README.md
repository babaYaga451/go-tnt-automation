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

# ⚙️ Back Pressure Handling in `go-tnt-tester`

This project is designed to efficiently process thousands of transit test records using Go pipelines. A key feature is how it **naturally handles back pressure** using **buffered channels and worker goroutines**.

---

## 🔄 Pipeline Architecture

Your pipeline stages are:

Each stage communicates through a **bounded buffered channel**:

```go
fileCh   := make(chan string, 1000)
sampleCh := make(chan Record, 10000)
enrichCh := make(chan Record, 10000)
resultCh := make(chan TestResult, 64000)
```

# ⚙️ Back Pressure Handling in go-tnt-automation

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
