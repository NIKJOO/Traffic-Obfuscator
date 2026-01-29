package main

import (
    "context"
    "fmt"
    "io"
    "math/rand"
    "net/http"
    "time"
    "sync"
    "github.com/chromedp/chromedp"
)

type Config struct {
    HeavyFiles []string
    EnableHeadlessBrowser bool
    ProxyAddress string
    BrowserTestDuration time.Duration
}

type TrafficSimulator struct {
    config     *Config
    httpClient *http.Client
    mode       string
    mu         sync.RWMutex
}

func NewTrafficSimulator(config *Config) *TrafficSimulator {
    transport := &http.Transport{
        Proxy: http.ProxyFromEnvironment,
    }

    return &TrafficSimulator{
        config: config,
        httpClient: &http.Client{
            Transport: transport,
            Timeout:   60 * time.Second,
        },
        mode: "idle",
    }
}

func (ts *TrafficSimulator) StartCamouflage() {
    ts.mu.Lock()
    defer ts.mu.Unlock()
    if ts.mode == "camouflage" {
        return
    }
    ts.mode = "camouflage"
    fmt.Println("ACTIVATING HYBRID MODE (Headless Browser + HTTP)")

    for i := 0; i < 3; i++ {
        go ts.fileDownloadWorker(i)
    }

    if ts.config.EnableHeadlessBrowser {
        go ts.browserWorker()
    }
}

func (ts *TrafficSimulator) StopCamouflage() {
    ts.mu.Lock()
    defer ts.mu.Unlock()
    ts.mode = "idle"
    fmt.Println("[~] Camouflage stopped.")
}

func (ts *TrafficSimulator) fileDownloadWorker(id int) {
    for {
        ts.mu.RLock()
        stop := ts.mode != "camouflage"
        ts.mu.RUnlock()
        if stop {
            return
        }

        targetURL := ts.config.HeavyFiles[rand.Intn(len(ts.config.HeavyFiles))]
        ts.downloadHeavyChunk(targetURL, id)
        
        time.Sleep(time.Duration(1+rand.Intn(2)) * time.Second)
    }
}

func (ts *TrafficSimulator) downloadHeavyChunk(urlStr string, workerID int) {
    chunkSize := int64(20 * 1024 * 1024) // 20 MB
    
    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

    start := time.Now()
    resp, err := ts.httpClient.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    n, _ := io.CopyN(io.Discard, resp.Body, chunkSize)
    duration := time.Since(start)
    mbps := (float64(n) * 8 / 1024 / 1024) / duration.Seconds()

    fmt.Printf("[HTTP Worker %d] Downloaded %d MB (%.2f Mbps)\n", workerID, n/1024/1024, mbps)
}

func (ts *TrafficSimulator) browserWorker() {
    for {
        ts.mu.RLock()
        stop := ts.mode != "camouflage"
        ts.mu.RUnlock()
        if stop {
            return
        }

        fmt.Println("[BROWSER] Launching Headless Chrome for fast.com...")
        speed, err := runFastComTest(ts.config)
        if err != nil {
            fmt.Printf("[BROWSER] Error: %v\n", err)
        } else {
            fmt.Printf("[BROWSER] Test finished. Detected Speed: %s Mbps\n", speed)
        }

        time.Sleep(10 * time.Second)
    }
}


func runFastComTest(config *Config) (string, error) {
    opts := []chromedp.ExecAllocatorOption{
        chromedp.NoFirstRun,
        chromedp.NoDefaultBrowserCheck,
        chromedp.Headless, 
        chromedp.DisableGPU,
        chromedp.NoSandbox,
        chromedp.IgnoreCertErrors,
    }

    if config.ProxyAddress != "" {
        opts = append(opts, chromedp.ProxyServer(config.ProxyAddress))
    }

    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()

    ctx, cancel := chromedp.NewContext(allocCtx)
    defer cancel()

    ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
    defer cancel()

    var speedText string

    err := chromedp.Run(ctx,
        chromedp.Navigate("https://fast.com"),
        chromedp.WaitVisible(`span.succeeded-color`, chromedp.ByQuery),
        chromedp.Sleep(config.BrowserTestDuration),
        chromedp.Text(`span.succeeded-color`, &speedText, chromedp.ByQuery),
    )

    return speedText, err
}

func main() {
    config := &Config{
        HeavyFiles: []string{
            "https://proof.ovh.net/files/100Mb.dat",
            "https://speed.hetzner.de/100MB.bin",
			"https://speedtest.ftp.otenet.gr/files/test100Mb.db",
        },
        EnableHeadlessBrowser: true,
        ProxyAddress:          "", 
        BrowserTestDuration:   15 * time.Second, 
    }

    

    fmt.Println("Press ENTER to Start Decoy ...")
    fmt.Scanln()
	
	sim := NewTrafficSimulator(config)
    sim.StartCamouflage()

    fmt.Println("Press ENTER to Stop...")
    fmt.Scanln()
    sim.StopCamouflage()
}
