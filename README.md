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

📚 [See Architecture Overview in the Wiki →](https://github.com/babaYaga451/go-tnt-automation.wiki.git)

![Visit Wiki](https://img.shields.io/badge/docs-wiki-blue?logo=github)
