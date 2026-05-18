package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost                 string
	DBPort                 string
	DBUser                 string
	DBPassword             string
	DBName                 string
	KafkaBroker            string
	KafkaTopicOrderCreated string
}

func LoadConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		DBHost:                 os.Getenv("DB_HOST"),
		DBPort:                 os.Getenv("DB_PORT"),
		DBUser:                 os.Getenv("DB_USER"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBName:                 os.Getenv("DB_NAME"),
		KafkaBroker:            os.Getenv("KAFKA_BROKER"),
		KafkaTopicOrderCreated: os.Getenv("KAFKA_TOPIC_ORDER_CREATED"),
	}
}
