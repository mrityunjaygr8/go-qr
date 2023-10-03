package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error reading configuration", err)
	}

	config := DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Password: os.Getenv("DB_PASS"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USER"),
	}

	store, err := NewPostgresStore(config)

	if err != nil {
		log.Fatal("error connecting to database", err)
	}
	server := &PlayerServer{store}
	port := "5000"
	log.Printf("Server starting on %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), server))
}
