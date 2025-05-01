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
