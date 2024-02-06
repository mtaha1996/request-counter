package config

type Config struct {
	WindowSizeInSeconds int
	PersistInterval     int
	DataTTL             int
	StoragePath         string
}

func (cfg Config) DefaultConfig() Config {
	return Config{
		WindowSizeInSeconds: 60,
		PersistInterval:     1,
		DataTTL:             60,
		StoragePath:         "storage/storage.json",
	}
}
