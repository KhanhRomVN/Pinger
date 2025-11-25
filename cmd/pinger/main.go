package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/pinger/internal/config"
	"github.com/yourusername/pinger/internal/logger"
	"github.com/yourusername/pinger/internal/pinger"
	"github.com/yourusername/pinger/internal/server"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer log.Sync()

	log.Info("Starting Pinger service",
		zap.Int("target_count", len(cfg.PingURLs)),
		zap.Duration("interval", cfg.PingInterval),
	)

	// Create pinger
	p := pinger.New(
		cfg.PingURLs,
		cfg.PingInterval,
		cfg.RequestTimeout,
		cfg.MaxRetries,
		cfg.LogResponseBody,
		log,
	)

	// Create HTTP server
	srv := server.New(cfg.Port, log)

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info("Received shutdown signal")
		cancel()
		
		// Shutdown HTTP server
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("HTTP server shutdown error", zap.Error(err))
		}
	}()

	// Start HTTP server in goroutine
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server error", zap.Error(err))
			cancel()
		}
	}()

	// Start pinger
	p.Start(ctx)

	log.Info("Pinger service stopped gracefully")
	return nil
}