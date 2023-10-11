package main

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
)

//func TestMain(m *testing.M) {
//	code, err := run(m)
//	if err != nil {
//		fmt.Println(err)
//	}
//	os.Exit(code)
//}

var db *sql.DB

const initialMigration = `
CREATE TABLE Scores (
	Name varchar UNIQUE NOT NULL ,
	Score int default 1 NOT NULL
);
`

const teardownMigration = `
DROP TABLE Scores;
`

func run(m *testing.M) (code int, err error) {
	dbConfig := DbConfig{
		Host:     "localhost",
		Password: "postgres",
		Port:     "15432",
		Name:     "postgres",
		Username: "postgres",
	}

	db, err = sql.Open("postgres", dbConfig.getDBString())
	if err != nil {
		return -1, fmt.Errorf("could not connect to database: %w", err)
	}

	_, err = db.Exec(initialMigration)
	if err != nil {
		log.Println("yo, here we are")
		//_, _ = db.Exec(teardownMigration)
		//
		//db.Close()
		return -1, fmt.Errorf("error creating tables: %w", err)
	}

	// truncates all test data after the tests are run
	defer func() {
		_, _ = db.Exec(teardownMigration)

		db.Close()
	}()

	return m.Run(), nil

}
