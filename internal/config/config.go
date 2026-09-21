package config

import (
	"fmt"
	"os"
	"sync"

	_ "github.com/joho/godotenv/autoload"
)

var (
	loaded Config
	once   sync.Once
)

type Config struct {
	Debug         bool
	JsonLogger    bool
	ReadOnly      bool
	ServerAddress string
	DBString      string
}

func Get() Config {
	return MustLoad()
}

func MustLoad() Config {
	once.Do(func() {
		loaded = load()
	})
	return loaded
}

func load() Config {
	return Config{
		Debug:         getenv("DEBUG", false) == "1",
		JsonLogger:    getenv("JSON_LOGGER", false) == "1",
		ReadOnly:      getenv("READ_ONLY", false) == "1",
		ServerAddress: getenv("SERVER_ADDRESS", true),
		DBString:      getenv("DB_STRING", true),
	}
}

func getenv(key string, required bool) string {
	value := os.Getenv(key)
	if value == "" && required {
		panic(fmt.Sprintf("config: %s variable is required but not set", key))
	}
	return value
}
