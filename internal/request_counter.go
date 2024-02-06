package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	config "github.com/request-counter/configs"
)

// // Constants for configuration
// const (
// 	WindowSizeInSeconds = 60 // Size of the moving window in seconds
// 	PersistInterval     = 1  // Interval for persisting data in seconds
// 	DataTTL             = 60 // Time-to-live for data in seconds
// 	cfg.storagePath         = "storage/storage.json"
// )

// ResponseJSON defines the structure of the response JSON
type ResponseJSON struct {
	Count int64 `json:"count"`
}

// StorageSchema represents the structure of data saved to the file
type StorageSchema struct {
	Data []map[int64]int64 `json:"data"`
}

// RequestCounter holds the request count data
type RequestCounter struct {
	data     sync.Map
	cfg      config.Config
	writeMux sync.RWMutex
}

// writeJSONToFile writes the provided data to a JSON file
func (rc *RequestCounter) writeJSONToFile(data any) error {
	file, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(rc.cfg.StoragePath, file, 0644)
}

// persist periodically saves the data to a file
func (rc *RequestCounter) PersistIntervaly() {
	ticker := time.NewTicker(time.Duration(rc.cfg.PersistInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rc.PersistOnFile()
	}
}

func (rc *RequestCounter) PersistOnFile() {

	log.Println("Persisting on file...")

	var storage StorageSchema
	rc.data.Range(func(key, value interface{}) bool {
		storage.Data = append(storage.Data, map[int64]int64{key.(int64): value.(int64)})
		return true
	})

	if err := rc.writeJSONToFile(storage); err != nil {
		log.Println("Error writing to file:", err)
	}

}

// fetch60SecCount returns the count of requests in the last 60 seconds
func (rc *RequestCounter) fetch60SecCount(now time.Time) int64 {
	var count int64
	rc.data.Range(func(key, value interface{}) bool {
		if key.(int64) > int64(now.Add(-time.Second*time.Duration(rc.cfg.WindowSizeInSeconds)).Unix()) &&
			key.(int64) <= int64(now.Unix()) {
			count += value.(int64)
		}
		return true
	})
	return count
}

// store increments the request count for the current second
func (rc *RequestCounter) store(now time.Time) {
	timestamp := now.Unix()
	rc.writeMux.Lock()
	defer rc.writeMux.Unlock()
	if count, exist := rc.data.LoadOrStore(timestamp, int64(1)); exist {
		rc.data.Store(timestamp, count.(int64)+1)
	}
}

// expiredRemover periodically removes expired data
func (rc *RequestCounter) ExpiredRemover() {
	ticker := time.NewTicker(time.Duration(rc.cfg.DataTTL) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rc.data.Range(func(key, value interface{}) bool {
			now := time.Now()
			if key.(int64) < int64(now.Add(-time.Second*time.Duration(rc.cfg.WindowSizeInSeconds)).Unix()) {
				rc.data.Delete(key)
			}
			return true
		})
	}
}

// count handles the HTTP request and responds with the request count
func (rc *RequestCounter) Count(w http.ResponseWriter, req *http.Request) {
	now := time.Now()

	rc.store(now)
	cnt := rc.fetch60SecCount(now)

	resp := ResponseJSON{Count: cnt}
	res, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Write(res)
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
func LoadStorage(cfg config.Config) *RequestCounter {
	rc := RequestCounter{cfg: cfg}

	storageFile, err := os.ReadFile(cfg.StoragePath)
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
