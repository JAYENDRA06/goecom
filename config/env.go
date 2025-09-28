package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	PublicHost              string
	Port                    string
	DBUser                  string
	DBPassword              string
	DBAddress               string
	DBName                  string
	JWTExpirationsInSeconds int64
	JWTSecret               string
}

var Envs *Config

func InitConfig() Config {
	fmt.Println("Initializing env variables")
	if Envs == nil {
		Envs = &Config{
			PublicHost:              getEnv("PUBLIC_HOST", "http://localhost"),
			Port:                    getEnv("PORT", "8080"),
			DBUser:                  getEnv("DB_USER", "root"),
			DBPassword:              getEnv("DB_PASSWORD", "Jayendra06*"),
			DBAddress:               fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
			DBName:                  getEnv("DB_NAME", "ecom"),
			JWTExpirationsInSeconds: getEnvsInt("JWT_EXP", 3600*24*7),
			JWTSecret:               getEnv("JWT_SECRET", "jayendra"),
		}
	}
	return *Envs
}

func getEnv(key, fallBack string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	log.Printf("Using fallback for %s", key)
	return fallBack
}

func getEnvsInt(key string, fallBack int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Printf("Failed parsing so using fallback for %s", key)
			return fallBack
		}
		return i
	}
	log.Printf("Using fallback for %s", key)
	return fallBack
}
