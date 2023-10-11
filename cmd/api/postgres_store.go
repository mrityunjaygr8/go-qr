package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var ErrPostgresConfigInvalid = errors.New("the provided postgres config is invalid")

type DbConfig struct {
	Host     string
	Password string
	Port     string
	Name     string
	Username string
}

func (d DbConfig) getDBString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", d.Username, d.Password, d.Host, d.Port, d.Name)
}

func (d DbConfig) validate() error {
	if d.Port == "" || d.Host == "" {
		return ErrPostgresConfigInvalid
	}

	return nil
}

func NewPostgresStore(dbConfig DbConfig) (*PostgresStore, error) {
	err := dbConfig.validate()
	if err != nil {
		return nil, err
	}
	dbString := dbConfig.getDBString()
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{db}, nil
}

type PostgresStore struct {
	DB *sql.DB
}

func (p *PostgresStore) GetLeague() League {
	res, err := p.DB.Query(fmt.Sprint(`SELECT Name, Score FROM Scores`))
	if err != nil {
		log.Fatal(err)
	}

	defer res.Close()

	var league []Player
	for res.Next() {
		var (
			name  string
			score int
		)

		if err := res.Scan(&name, &score); err != nil {
			log.Fatal(err)
		}

		league = append(league, Player{name, score})
	}

	return league
}

func (p *PostgresStore) GetPlayerScore(playerName string) int {
	res, err := p.DB.Query(fmt.Sprintf(`SELECT Name, Score FROM Scores WHERE name = '%s'`, playerName))
	if err != nil {
		log.Fatal(err)
	}

	defer res.Close()

	var (
		name  string
		score int
	)
	for res.Next() {

		if err := res.Scan(&name, &score); err != nil {
			log.Fatal(err)
		}
	}

	return score
}

func (p *PostgresStore) RecordWin(name string) {
	score := p.GetPlayerScore(name)

	if score == 0 {
		_, err := p.DB.Exec(fmt.Sprintf(`INSERT INTO Scores (Name, Score) VALUES ('%s', %d)`, name, score+1))
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	_, err := p.DB.Exec(fmt.Sprintf(`UPDATE Scores SET Score = %d WHERE Name = '%s'`, score+1, name))
	if err != nil {
		log.Fatal(err)
	}
	return

}
