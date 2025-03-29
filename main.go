package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/2k4sm/remoteDictionary/src/api"
	"github.com/2k4sm/remoteDictionary/src/config"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg := config.LoadConfig()

	log.Printf("Starting remoteDictionary with port=%s, maxCacheSize=%d, maxKeySize=%d, maxValueSize=%d",
		cfg.Port, cfg.MaxCacheSize, cfg.MaxKeySize, cfg.MaxValueSize)

	server := api.NewServer(cfg)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
