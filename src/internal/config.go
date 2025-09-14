package internal

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	MongoURI       string
	RedisURI       string
	JWTSecret      string
	TxtBeltAPIKey  string
	OutputFileName string
}

var (
	once sync.Once
	c    *Config
)

func GetConfig() *Config {
	once.Do(func() {
		_ = godotenv.Load()
		c = &Config{
			Port:           must("PORT"),
			MongoURI:       must("MONGO_URI"),
			JWTSecret:      must("JWT_SECRET"),
			RedisURI:       must("REDIS_URI"),
			TxtBeltAPIKey:  must("TXT_BELT_API_KEY"),
			OutputFileName: must("OUTPUT_FILE_NAME"),
		}
	})
	if c == nil {
		log.Panic("config not loaded")
	}
	return c
}
func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("missing required env: %s", k)
	}
	return v
}
