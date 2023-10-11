package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const dbName = "game.db.json"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error reading configuration", err)
	}

	//config := DbConfig{
	//	Host:     os.Getenv("DB_HOST"),
	//	Password: os.Getenv("DB_PASS"),
	//	Port:     os.Getenv("DB_PORT"),
	//	Name:     os.Getenv("DB_NAME"),
	//	Username: os.Getenv("DB_USER"),
	//}
	//
	//store, err := NewPostgresStore(config)
	//
	//if err != nil {
	//	log.Fatal("error connecting to database", err)
	//}

	db, err := os.OpenFile(dbName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("error opening %s %v", dbName, err)
	}

	store, err := NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}
	server := NewPlayerServer(store)
	port := "5000"
	log.Printf("Server starting on %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), server))
}
