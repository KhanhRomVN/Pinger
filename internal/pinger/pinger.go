package pinger

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Pinger struct {
	urls           []string
	interval       time.Duration
	timeout        time.Duration
	maxRetries     int
	logResponseBody bool
	logger         *zap.Logger
	client         *http.Client
}

type PingResult struct {
	URL          string
	StatusCode   int
	ResponseTime time.Duration
	Success      bool
	Error        error
	Body         string
}

func New(urls []string, interval, timeout time.Duration, maxRetries int, logResponseBody bool, logger *zap.Logger) *Pinger {
	return &Pinger{
		urls:            urls,
		interval:        interval,
		timeout:         timeout,
		maxRetries:      maxRetries,
		logResponseBody: logResponseBody,
		logger:          logger,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (p *Pinger) Start(ctx context.Context) {
	p.logger.Info("Pinger started",
		zap.Int("url_count", len(p.urls)),
		zap.Duration("interval", p.interval),
		zap.Duration("timeout", p.timeout),
	)

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	// Ping ngay lập tức khi start
	p.pingAll()

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Pinger stopped")
			return
		case <-ticker.C:
			p.pingAll()
		}
	}
}

func (p *Pinger) pingAll() {
	var wg sync.WaitGroup
	results := make(chan PingResult, len(p.urls))

	for _, url := range p.urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			result := p.pingURL(u)
			results <- result
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		p.logResult(result)
	}
}

func (p *Pinger) pingURL(url string) PingResult {
	var lastErr error
	
	for attempt := 0; attempt <= p.maxRetries; attempt++ {
		if attempt > 0 {
			p.logger.Debug("Retrying request",
				zap.String("url", url),
				zap.Int("attempt", attempt),
			)
			time.Sleep(time.Second * time.Duration(attempt))
		}

		start := time.Now()
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		req.Header.Set("User-Agent", "Pinger/1.0")

		resp, err := p.client.Do(req)
		responseTime := time.Since(start)

		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		defer resp.Body.Close()

		var body string
		if p.logResponseBody {
			bodyBytes, _ := io.ReadAll(resp.Body)
			body = string(bodyBytes)
		}

		success := resp.StatusCode >= 200 && resp.StatusCode < 300

		return PingResult{
			URL:          url,
			StatusCode:   resp.StatusCode,
			ResponseTime: responseTime,
			Success:      success,
			Error:        nil,
			Body:         body,
		}
	}

	return PingResult{
		URL:     url,
		Success: false,
		Error:   lastErr,
	}
}

func (p *Pinger) logResult(result PingResult) {
	fields := []zap.Field{
		zap.String("url", result.URL),
		zap.Bool("success", result.Success),
	}

	if result.Success {
		fields = append(fields,
			zap.Int("status_code", result.StatusCode),
			zap.Duration("response_time", result.ResponseTime),
		)
		
		if p.logResponseBody && result.Body != "" {
			fields = append(fields, zap.String("body", result.Body))
		}

		p.logger.Info("Ping successful", fields...)
	} else {
		if result.Error != nil {
			fields = append(fields, zap.Error(result.Error))
		}
		if result.StatusCode > 0 {
			fields = append(fields,
				zap.Int("status_code", result.StatusCode),
				zap.Duration("response_time", result.ResponseTime),
			)
		}
		p.logger.Error("Ping failed", fields...)
	}
}