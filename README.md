# Pinger 🏓

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/KhanhRomVN/Pinger)](https://goreportcard.com/report/github.com/KhanhRomVN/Pinger)

> Lightweight, concurrent HTTP health checker built with Go. Keep your services alive with configurable intervals, retry logic, and structured logging.

![Pinger Demo](https://via.placeholder.com/1200x400/00ADD8/ffffff?text=Pinger+-+HTTP+Health+Monitoring)

---

## 📖 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Quick Start](#-quick-start)
- [Configuration](#-configuration)
- [Usage](#-usage)
- [Architecture](#-architecture)
- [Deployment](#-deployment)
- [Development](#-development)
- [Troubleshooting](#-troubleshooting)
- [License](#-license)
- [Contact](#-contact)

---

## 🌟 Overview

**Pinger** is a production-ready HTTP health monitoring service written in Go. It periodically sends GET requests to configured endpoints to ensure they remain responsive. Perfect for keeping free-tier services (Render, Heroku, Railway) awake or monitoring critical APIs.

### Why Pinger?

- **Lightweight**: Single binary, minimal resource usage (~5MB memory)
- **Concurrent**: Parallel health checks with goroutines
- **Resilient**: Configurable retry logic with exponential backoff
- **Observable**: Structured logging with zap (JSON or console)
- **Production-Ready**: Graceful shutdown, signal handling, error recovery
- **Easy Deploy**: Docker support, systemd service, cloud-ready

---

## ✨ Features

### 🚀 Core Functionality
- **Multi-Target Monitoring**: Check unlimited URLs concurrently
- **Configurable Intervals**: Set custom ping frequencies (seconds/minutes/hours)
- **Smart Retries**: Automatic retry with exponential backoff
- **Request Timeout**: Prevent hanging requests with configurable timeout
- **Response Logging**: Optional body logging for debugging

### 📊 Monitoring & Logging
- **Structured Logs**: JSON or console output with zap
- **Log Levels**: Debug, Info, Warn, Error
- **Rich Context**: Status codes, response times, error details
- **Color Output**: Terminal-friendly colored logs

### 🛡️ Reliability
- **Graceful Shutdown**: SIGTERM/SIGINT handling
- **Context Cancellation**: Clean goroutine cleanup
- **Error Recovery**: Continues on individual failures
- **Concurrent-Safe**: Thread-safe operations

### ⚙️ Configuration
- **Environment Variables**: 12-factor app compliant
- **.env File Support**: Easy local development
- **Validation**: Config validation at startup
- **Flexible**: Override defaults per environment

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21+ (for building from source)
- Or Docker (for containerized deployment)

### Installation

#### Option 1: Download Binary (Recommended)

```bash
# Download latest release (replace VERSION)
wget https://github.com/KhanhRomVN/Pinger/releases/download/v1.0.0/pinger-linux-amd64

# Make executable
chmod +x pinger-linux-amd64
mv pinger-linux-amd64 /usr/local/bin/pinger
```

#### Option 2: Build from Source

```bash
# Clone repository
git clone https://github.com/KhanhRomVN/Pinger.git
cd Pinger

# Install dependencies
go mod download

# Build binary
go build -o pinger cmd/pinger/main.go
```

#### Option 3: Docker

```bash
# Pull image
docker pull khanhromvn/pinger:latest

# Run container
docker run -d \
  -e PING_URLS="https://api.example.com/health" \
  -e PING_INTERVAL=60 \
  khanhromvn/pinger:latest
```

### Basic Usage

```bash
# Create .env file
cp .env.example .env

# Edit configuration
nano .env

# Run pinger
./pinger
```

**Output:**
```
2024-01-15T10:30:00.123Z  INFO  Starting Pinger service  target_count=3  interval=1m0s
2024-01-15T10:30:00.456Z  INFO  Pinger started  url_count=3  interval=1m0s  timeout=10s
2024-01-15T10:30:01.234Z  INFO  Ping successful  url=https://api.example.com/health  success=true  status_code=200  response_time=778ms
```

---

## ⚙️ Configuration

### Environment Variables

Create a `.env` file or set environment variables:

```bash
# Required: Comma-separated list of URLs
PING_URLS=https://api.example.com/health,https://service2.com/ping

# Optional: Ping interval in seconds (default: 60)
PING_INTERVAL=60

# Optional: Request timeout in seconds (default: 10)
REQUEST_TIMEOUT=10

# Optional: Max retry attempts (default: 3)
MAX_RETRIES=3

# Optional: Log level: debug|info|warn|error (default: info)
LOG_LEVEL=info

# Optional: Log response body (default: false)
LOG_RESPONSE_BODY=false
```

### Configuration Details

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `PING_URLS` | string | *required* | Comma-separated URLs to monitor |
| `PING_INTERVAL` | int | 60 | Seconds between each ping cycle |
| `REQUEST_TIMEOUT` | int | 10 | HTTP request timeout in seconds |
| `MAX_RETRIES` | int | 3 | Retry attempts before marking as failed |
| `LOG_LEVEL` | string | info | Logging verbosity level |
| `LOG_RESPONSE_BODY` | bool | false | Include response body in logs |

### Example Configurations

**Keep Free Services Alive:**
```bash
# Ping every 5 minutes to prevent sleep
PING_URLS=https://myapp.onrender.com/health
PING_INTERVAL=300
REQUEST_TIMEOUT=30
MAX_RETRIES=5
```

**API Monitoring:**
```bash
# Monitor multiple APIs with detailed logging
PING_URLS=https://api1.com/v1/health,https://api2.com/status,https://api3.com/ping
PING_INTERVAL=30
LOG_LEVEL=debug
LOG_RESPONSE_BODY=true
```

**Production Setup:**
```bash
# Production monitoring with quick failure detection
PING_URLS=https://prod-api.com/health
PING_INTERVAL=10
REQUEST_TIMEOUT=5
MAX_RETRIES=2
LOG_LEVEL=error
```

---

## 📚 Usage

### Running Locally

```bash
# Simple run
./pinger

# With environment variables
PING_URLS=https://example.com LOG_LEVEL=debug ./pinger

# Run in background
nohup ./pinger > pinger.log 2>&1 &
```

### Systemd Service

Create `/etc/systemd/system/pinger.service`:

```ini
[Unit]
Description=Pinger Health Monitoring Service
After=network.target

[Service]
Type=simple
User=pinger
WorkingDirectory=/opt/pinger
EnvironmentFile=/opt/pinger/.env
ExecStart=/opt/pinger/pinger
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Manage service:**
```bash
# Enable and start
sudo systemctl enable pinger
sudo systemctl start pinger

# Check status
sudo systemctl status pinger

# View logs
sudo journalctl -u pinger -f
```

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  pinger:
    image: khanhromvn/pinger:latest
    container_name: pinger
    environment:
      - PING_URLS=https://api.example.com/health,https://service2.com/ping
      - PING_INTERVAL=60
      - REQUEST_TIMEOUT=10
      - MAX_RETRIES=3
      - LOG_LEVEL=info
    restart: unless-stopped
```

**Run:**
```bash
docker-compose up -d
```

---

## 🏗️ Architecture

### Project Structure

```
Pinger/
├── cmd/
│   └── pinger/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration loader
│   ├── logger/
│   │   └── logger.go         # Zap logger setup
│   └── pinger/
│       └── pinger.go         # Core ping logic
├── .env.example              # Example configuration
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── Dockerfile                # Container build
└── README.md
```

### Component Overview

```
┌─────────────────────────────────────────────────────────┐
│                    Pinger Service                        │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────────┐      ┌──────────────┐      ┌────────┐ │
│  │    Main      │─────►│    Config    │─────►│ Logger │ │
│  │  (cmd/main)  │      │   Loader     │      │  (zap) │ │
│  └──────────────┘      └──────────────┘      └────────┘ │
│         │                                                 │
│         │ creates                                         │
│         ▼                                                 │
│  ┌──────────────────────────────────────────────┐       │
│  │           Pinger Core                         │       │
│  │  - Start(ctx) : goroutine per URL            │       │
│  │  - pingAll()  : concurrent execution         │       │
│  │  - pingURL()  : retry logic + timeout        │       │
│  │  - logResult(): structured logging           │       │
│  └──────────────────────────────────────────────┘       │
│         │                                                 │
│         │ HTTP GET                                        │
│         ▼                                                 │
└─────────────────────────────────────────────────────────┘
          │
          │ Requests
          ▼
   ┌──────────────────┐
   │  Target Services │
   │  (URLs to ping)  │
   └──────────────────┘
```

### Flow Diagram

```
Start
  │
  ├─ Load Config (.env)
  │   └─ Validate PING_URLS
  │
  ├─ Initialize Logger (zap)
  │   └─ Set log level
  │
  ├─ Create Pinger Instance
  │   ├─ HTTP Client (with timeout)
  │   └─ Ticker (interval)
  │
  ├─ Setup Context + Signal Handler
  │   └─ Listen for SIGTERM/SIGINT
  │
  ├─ Start Pinger
  │   ├─ Immediate ping (t=0)
  │   └─ Ticker loop
  │       ├─ Every PING_INTERVAL
  │       │   └─ pingAll()
  │       │       └─ Spawn goroutine per URL
  │       │           ├─ Retry loop (MAX_RETRIES)
  │       │           │   ├─ HTTP GET
  │       │           │   ├─ Check status code
  │       │           │   └─ Measure response time
  │       │           └─ Log result
  │       │
  │       └─ Wait for signal
  │           └─ Graceful shutdown
  │
  └─ Exit
```

---

## 🌐 Deployment

### Deploy to Render

**render.yaml:**
```yaml
services:
  - type: web
    name: pinger
    env: docker
    dockerfilePath: ./Dockerfile
    envVars:
      - key: PING_URLS
        value: https://myapp.onrender.com/health
      - key: PING_INTERVAL
        value: 300
      - key: LOG_LEVEL
        value: info
```

### Deploy to Railway

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login
railway login

# Deploy
railway up
```

### Deploy to Heroku

```bash
# Create app
heroku create my-pinger

# Set config
heroku config:set PING_URLS=https://example.com/health

# Deploy
git push heroku main
```

### VPS Deployment

```bash
# Upload binary
scp pinger user@server:/opt/pinger/

# Upload .env
scp .env user@server:/opt/pinger/

# Setup systemd service (see Usage section)
sudo systemctl enable pinger
sudo systemctl start pinger
```

### Kubernetes Deployment

**deployment.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pinger
spec:
  replicas: 1
  selector:
    matchLabels:
      app: pinger
  template:
    metadata:
      labels:
        app: pinger
    spec:
      containers:
      - name: pinger
        image: khanhromvn/pinger:latest
        env:
        - name: PING_URLS
          value: "https://api.example.com/health"
        - name: PING_INTERVAL
          value: "60"
        resources:
          limits:
            memory: "32Mi"
            cpu: "50m"
```

---

## 🛠️ Development

### Building

```bash
# Build for current platform
go build -o pinger cmd/pinger/main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o pinger-linux-amd64 cmd/pinger/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o pinger-windows-amd64.exe cmd/pinger/main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o pinger-darwin-amd64 cmd/pinger/main.go
```

### Docker Build

```bash
# Build image
docker build -t khanhromvn/pinger:latest .

# Run locally
docker run --rm \
  -e PING_URLS=https://example.com \
  khanhromvn/pinger:latest

# Push to registry
docker push khanhromvn/pinger:latest
```

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Benchmark
go test -bench=. ./internal/pinger
```

### Code Quality

```bash
# Format code
go fmt ./...

# Lint
golangci-lint run

# Vet
go vet ./...
```

---

## 🐛 Troubleshooting

### Pinger Won't Start

**Problem:** `.env file not found` error

**Solution:**
```bash
# Create .env from example
cp .env.example .env

# Or set environment variables directly
export PING_URLS=https://example.com
./pinger
```

### Connection Timeout

**Problem:** `request failed: context deadline exceeded`

**Solution:**
```bash
# Increase timeout
REQUEST_TIMEOUT=30
MAX_RETRIES=5
```

### High Memory Usage

**Problem:** Memory usage increasing over time

**Solution:**
```bash
# Disable response body logging
LOG_RESPONSE_BODY=false

# Reduce log level
LOG_LEVEL=error
```

### Permission Denied (Systemd)

**Problem:** `Permission denied` when accessing logs

**Solution:**
```bash
# Create dedicated user
sudo useradd -r -s /bin/false pinger

# Fix permissions
sudo chown -R pinger:pinger /opt/pinger
sudo chmod 755 /opt/pinger/pinger
```

### Ping Fails But Service is Up

**Problem:** Ping reports failure but manual curl works

**Solution:**
```bash
# Check User-Agent blocking
# Pinger uses "Pinger/1.0" by default

# Test with curl
curl -H "User-Agent: Pinger/1.0" https://example.com

# If blocked, modify user agent in code (pinger.go:85)
```

---

## 📊 Use Cases

### 1. Keep Free Tier Services Alive
```bash
# Prevent Render/Heroku free tier sleep
PING_URLS=https://myapp.onrender.com/health
PING_INTERVAL=300  # Every 5 minutes
```

### 2. API Health Monitoring
```bash
# Monitor critical APIs
PING_URLS=https://api1.com/v1/health,https://api2.com/status
PING_INTERVAL=30
LOG_LEVEL=info
```

### 3. Load Testing Preparation
```bash
# Warm up services before load test
PING_URLS=https://staging.example.com/api/warmup
PING_INTERVAL=5
REQUEST_TIMEOUT=60
```

### 4. Uptime Monitoring
```bash
# Simple uptime check
PING_URLS=https://example.com
PING_INTERVAL=60
LOG_RESPONSE_BODY=false
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2024 KhanhRomVN

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## 📧 Contact

**Khánh Rom**
- GitHub: [@KhanhRomVN](https://github.com/KhanhRomVN)
- Email: khanhromvn@gmail.com

**Project Link:** [https://github.com/KhanhRomVN/Pinger](https://github.com/KhanhRomVN/Pinger)

---

## 🙏 Acknowledgments

- [Uber Zap](https://github.com/uber-go/zap) - Blazing fast structured logging
- [godotenv](https://github.com/joho/godotenv) - .env file support
- [Go Standard Library](https://pkg.go.dev/std) - Built on solid foundations

---

## 🚀 Roadmap

- [ ] Prometheus metrics export
- [ ] Webhook notifications on failure
- [ ] HTTP POST support with custom payloads
- [ ] Response validation (status code, body pattern)
- [ ] Dashboard for visualization
- [ ] Multi-region pinging
- [ ] Slack/Discord integration

---

<div align="center">

Made with ❤️ by [KhanhRomVN](https://github.com/KhanhRomVN)

⭐ Star this repo if you find it helpful!

**Keep your services alive, one ping at a time!** 🏓

</div>