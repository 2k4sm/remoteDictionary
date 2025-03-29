package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	MaxKeySize   int
	MaxValueSize int
}

func LoadConfig() *Config {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "7171"
	}

	maxKeySize := 256
	if size := os.Getenv("MAX_KEY_SIZE"); size != "" {
		if val, err := strconv.Atoi(size); err == nil {
			if val <= 256 {
				maxKeySize = val
			} else {
				log.Printf("Warning: MAX_KEY_SIZE exceeds limit of 256, using 256")
			}
		}
	}

	maxValueSize := 256
	if size := os.Getenv("MAX_VALUE_SIZE"); size != "" {
		if val, err := strconv.Atoi(size); err == nil {
			if val <= 256 {
				maxValueSize = val
			} else {
				log.Printf("Warning: MAX_VALUE_SIZE exceeds limit of 256, using 256")
			}
		}
	}

	return &Config{
		Port:         port,
		MaxKeySize:   maxKeySize,
		MaxValueSize: maxValueSize,
	}
}
