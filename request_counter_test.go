package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

var Cfg = Config{}.DefaultConfig()

func TestRequestCounterStoreAndFetch(t *testing.T) {
	rc := RequestCounter{data: sync.Map{}, cfg: Cfg}

	// Simulate requests at different times
	now := time.Now()
	for i := 0; i < 10; i++ {
		rc.store(now.Add(-time.Duration(i) * time.Second))
	}

	count := rc.fetch60SecCount(now)
	if count != 10 {
		t.Errorf("Expected count to be 10, got %d", count)
	}

	// Test count outside the 60-second window
	oldCount := rc.fetch60SecCount(now.Add(-61 * time.Second))
	if oldCount != 0 {
		t.Errorf("Expected count to be 0 for old timestamp, got %d", oldCount)
	}
}

func TestPersistence(t *testing.T) {

	cfg := Cfg

	defer os.Remove(cfg.storagePath)
	os.Remove(cfg.storagePath)

	rc := RequestCounter{data: sync.Map{}, cfg: cfg}

	// Simulate some requests
	now := time.Now()
	rc.store(now)
	rc.store(now.Add(-30 * time.Second)) // 30 seconds earlier

	// Persist data to file
	go rc.Persist()
	time.Sleep(time.Duration(5*cfg.persistInterval) * time.Second)

	// Load data from file
	loadedRc := LoadStorage(cfg)

	count := loadedRc.fetch60SecCount(now)
	if count != 2 {
		t.Errorf("Expected count to be 2 after loading, got %d", count)
	}
}

func TestHTTPRequest(t *testing.T) {
	cfg := Cfg
	rc := RequestCounter{data: sync.Map{}, cfg: cfg}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rc.Count)

	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"count":1}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
