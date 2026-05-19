package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost                       string
	DBUser                       string
	DBPassword                   string
	DBName                       string
	DBPort                       string
	DBSchema                     string
	ChromedpHeadless             bool
	ChromedpDisableBlinkFeatures string
	ChromedpUserAgent            string
	RedisAddr                    string
	KafkaBrokers                 string
	KafkaTopic                   string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DBHost:                       getEnv("DB_HOST", "localhost"),
		DBUser:                       getEnv("DB_USER", "postgres"),
		DBPassword:                   getEnv("DB_PASSWORD", ""),
		DBName:                       getEnv("DB_NAME", "wb_db"),
		DBPort:                       getEnv("DB_PORT", "5432"),
		DBSchema:                     getEnv("DB_SCHEMA", "wb_scraper"),
		ChromedpHeadless:             getBoolEnv("CHROMEDP_HEADLESS", true),
		ChromedpDisableBlinkFeatures: getEnv("CHROMEDP_DISABLE_BLINK_FEATURES", "AutomationControlled"),
		ChromedpUserAgent:            getEnv("CHROMEDP_USER_AGENT", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		RedisAddr:                    getEnv("REDIS_ADDR", "localhost:6379"),
		KafkaBrokers:                 getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:                   getEnv("KAFKA_TOPIC", "wb-raw-products"),
	}
}

func getEnv(key, fallBack string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallBack
}

func getBoolEnv(key string, fallBack bool) bool {
	s := os.Getenv(key)
	if s == "" {
		return fallBack
	}

	value, err := strconv.ParseBool(s)
	if err != nil {
		return fallBack
	}
	return value
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort)
}
