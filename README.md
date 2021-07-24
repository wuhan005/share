# 🐑 share ![Go](https://github.com/wuhan005/share/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/wuhan005/share)](https://goreportcard.com/report/github.com/wuhan005/share) [![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?logo=sourcegraph)](https://sourcegraph.com/github.com/wuhan005/share)

Share files at once. / 懂得都懂

## Getting started

### Step 1. Install

```bash
go install
```

### Step 2. Set your account

Create these environment variables.

```bash
export CHAOXING_ACCOUNT_PHONE=<REDACTED>
export CHAOXING_ACCOUNT_PASSWORD=<REDACTED>
```

### Step 3. Upload a file

```bash
share Elaina-Houki.jpg
```

#### Output

```
 [TRACE] Upload file "Elaina-Houki.jpg"...
 100% |===========================================================================================| (4.8/4.8 MB, 3.995 MB/s)        
 [ INFO] http://d0.ananas.chaoxing.com/download/91247eea0d8b1fc6b48f1c2750a22b57?at_=1627155170160&ak_=5deafc4d55b2b1b127a0d5dfc422c1e5&ad_=111fd8770d9928bf2a1877868552025e
```

## License

MIT License