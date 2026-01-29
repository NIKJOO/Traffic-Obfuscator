# Traffic Obfuscator

A traffic simulation and obfuscation tool written in Go.  
It generates high-volume noise traffic to mask user behavior and bypass Deep Packet Inspection (DPI) systems that rely on behavioral analysis.

The tool uses a hybrid approach combining HTTP-based traffic generation with real browser automation to execute JavaScript and mimic human browsing behavior.

## Features

- Real browser simulation using chromedp with headless Chrome
- Executes real JavaScript (e.g. Fast.com speed tests)
- High-bandwidth download masking using OVH and Hetzner servers
- Hybrid traffic pattern combining heavy downloads and light browsing
- Automatic system proxy usage (e.g. V2rayN system proxy mode)
- Fully configurable worker counts, chunk sizes, and durations

## How It Works

### Problem

Uploading large files creates a distinct traffic pattern:

- High upload-to-download ratio
- Continuous traffic to a single destination

### Solution

The tool generates mixed traffic patterns:

#### HTTP Workers
- Download data from:
  - https://proof.ovh.net
  - https://speed.hetzner.de
- Consume downstream bandwidth
- Reduce effective upload ratio

#### Browser Worker
- Navigates to:
  - https://fast.com
  - https://speed.cloudflare.com
- Executes JavaScript speed tests
- Produces real browser-based traffic

The resulting traffic profile resembles a power user downloading content or performing network diagnostics.

## Prerequisites

- Go 1.19 or higher
- Google Chrome or Microsoft Edge
- Optional: System proxy software (e.g. V2rayN)

## Installation

```bash
git clone https://github.com/your-username/your-repo-name.git
cd your-repo-name
go get -u github.com/chromedp/chromedp
go build -o traffic-obfuscator main.go
```

### Usage

Run the application:
```bash
./traffic-obfuscator
```

Press ENTER to start traffic generation
Start your upload immediately after
Press ENTER again to stop traffic generation

### Technical Details

- HTTP Workers: Utilize net/http with io.CopyN to download specific chunks (e.g., 20MB) to simulate buffer downloading.
- Browser Worker: Uses Chrome DevTools Protocol (chromedp) to fetch pages, wait for DOM elements (like speed results), and extract data.
- Rate Limiting: Includes randomized delays to prevent detection as a bot/DDoS attack.

Randomized delays are used to avoid bot or DDoS detection

### Disclaimer

This tool is intended for educational and research purposes to demonstrate traffic obfuscation techniques. Use responsibly and in compliance with your local laws and network usage policies.
