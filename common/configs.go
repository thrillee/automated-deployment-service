package common

import "os"

type DbConfig struct {
	URI     string
	DB_NAME string
}

func GetDBConfig() *DbConfig {
	return &DbConfig{
		URI:     os.Getenv("DB_URI"),
		DB_NAME: os.Getenv("DB_NAME"),
	}
}
