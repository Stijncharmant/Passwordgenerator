package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBName     string `json:"db_name"`
}

type Password struct {
	Password string    `json:"password"`
	Date     time.Time `json:"date"`
}

func main() {
	// Load configuration from JSON file
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	// Open a connection to the database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}
	defer db.Close()

	// Generate a password
	password := generatePassword()

	// Store the password in the database
	err = storePassword(db, password)
	if err != nil {
		log.Fatalf("Failed to store password: %s", err)
	}

	fmt.Printf("Password generated: %s\n", password)
}

func loadConfig(filename string) (*Config, error) {
	config := &Config{}

	configFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open config file: %s", err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(config); err != nil {
		return nil, fmt.Errorf("Failed to decode config file: %s", err)
	}

	return config, nil
}

func generatePassword() string {
	const (
		length  = 12
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-="
	)

	password := make([]byte, length)
	rand.Read(password)

	for i, b := range password {
		password[i] = charset[int(b)%len(charset)]
	}

	return string(password)
}

func storePassword(db *sql.DB, password string) error {
	insertQuery := "INSERT INTO passwords (password, date) VALUES (?, ?)"
	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := db.Exec(insertQuery, password, now)
	if err != nil {
		return fmt.Errorf("Failed to execute insert query: %s", err)
	}

	return nil
}
