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

## Architecture

<details>
<summary>Click to expand Mermaid code</summary>
flowchart LR
    subgraph Pipeline Stages
        A[📁 discover\n(walk files)] --> B[🧪 sample\n(reservoir sampling)]
        B --> C[🔍 enrich\n(city/state/zip lookup)]
        C --> D[🌐 api call\n(HTTP request + validation)]
        D --> E[📝 result collect\n(JUnit writer)]
    end

    subgraph Channels
        a1[fileCh\n(chan string, 1000)]
        a2[sampleCh\n(chan Record, 10000)]
        a3[enrichCh\n(chan Record, 10000)]
        a4[resultCh\n(chan TestResult, 64000)]
    end

    A --> a1 --> B
    B --> a2 --> C
    C --> a3 --> D
    D --> a4 --> E

    %% Back pressure arrows
    style a4 stroke:#f00,stroke-width:2px
    style a3 stroke:#f00,stroke-width:2px
    style a2 stroke:#f00,stroke-width:2px
    style a1 stroke:#f00,stroke-width:2px

    D -. blocks on full resultCh .-> C
    C -. blocks on full enrichCh .-> B
    B -. blocks on full sampleCh .-> A
    A -. blocked write .-> STOP[🔁 pipeline slows down (back pressure)]

</details>
