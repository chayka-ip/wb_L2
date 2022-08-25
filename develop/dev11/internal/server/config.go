package server

//Config is a config for http server
type Config struct {
	Host           string `toml:"host"`
	Port           string `toml:"port"`
	CacheSizeBytes uint64 `toml:"cahe_size_bytes"`
	DiskvDir       string `toml:"diskv_dir"`
}

//NewDefaultConfig parses congif from config file
func NewDefaultConfig() *Config {
	conf := Config{
		Host:           "localhost",
		Port:           "8080",
		CacheSizeBytes: 1024 * 1024,
		DiskvDir:       "db_dir",
	}
	return &conf
}
