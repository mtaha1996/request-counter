// Package main implements an HTTP server that tracks the number of requests
// received in the last 60 seconds. It demonstrates concurrency handling,
// HTTP server implementation, and data persistence in Go.
// The server stores request counts in a synchronized map and persists this data
// to a file to maintain state across restarts.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ResponseJSON defines the structure of the response JSON
type ResponseJSON struct {
	Count int64 `json:"count"`
}

func convertStorageToSyncMap(storageSchema StorageSchema, rc *RequestCounter) *RequestCounter {

	for _, mapData := range storageSchema.Data {
		for key, val := range mapData {
			rc.data.Store(key, val)
		}
	}

	return rc
}

// LoadStorage loads the persisted data from the file
func LoadStorage(cfg Config) *RequestCounter {
	rc := RequestCounter{cfg: cfg}

	storageFile, err := os.ReadFile(cfg.storagePath)
	if err != nil {
		log.Println("Error opening storage file:", err)
		return &rc
	}

	var storageSchema StorageSchema
	if err = json.Unmarshal(storageFile, &storageSchema); err != nil {
		log.Println("Error parsing storage file:", err)
		return &rc
	}

	return convertStorageToSyncMap(storageSchema, &rc)
}

// main is the entry point of the application
func main() {

	// Setting up a channel to listen for termination signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	cfg := Config{}

	flag.IntVar(&cfg.windowSizeInSeconds, "window-size", 60, "Size of the moving window in seconds")
	flag.IntVar(&cfg.persistInterval, "persist-interval", 5, "Interval for persisting data in seconds")
	flag.IntVar(&cfg.dataTTL, "data-ttl", 60, "Time-to-live for data in seconds")
	flag.StringVar(&cfg.storagePath, "storage-path", "storage/storage.json", "Path to the storage file")

	flag.Parse()

	reqCounter := LoadStorage(cfg)

	go reqCounter.PersistIntervaly()
	go reqCounter.expiredRemover()

	serverAddress := ":8080"
	http.HandleFunc("/", reqCounter.Count)
	srv := &http.Server{Addr: serverAddress}

	// Start your HTTP server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Log that the server has started
	log.Printf("Server started successfully on http://localhost%s", serverAddress)

	// Block until a signal is received
	<-stopChan
	log.Println("Shutting down server...")
	reqCounter.PersistOnFile()

	// Initiate graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
