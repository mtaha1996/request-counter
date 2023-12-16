package main

type Config struct {
	windowSizeInSeconds int
	persistInterval     int
	dataTTL             int
	storagePath         string
}

func (cfg Config) DefaultConfig() Config {
	return Config{
		windowSizeInSeconds: 60,
		persistInterval:     1,
		dataTTL:             60,
		storagePath:         "storage/storage.json",
	}
}
