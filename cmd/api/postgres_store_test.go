package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostgresStore(t *testing.T) {
	store := PostgresStore{DB: db}
	t.Run("test initialization", func(t *testing.T) {
		err := store.DB.Ping()

		if err != nil {
			t.Errorf("tried db.ping, got %v", err)
		}

	})
}

func TestPostgresStore_GetPlayerScore(t *testing.T) {
	store := &PostgresStore{DB: db}
	server := NewPlayerServer(store)

	t.Run("get score of player in store", func(t *testing.T) {
		player := "Pepper"
		score := 5
		_, err := store.DB.Exec(fmt.Sprintf(`INSERT INTO Scores (Name, Score) VALUES ('%s', %d)`, player, score))
		require.Nil(t, err)

		request := newGetScoreRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		require.Equal(t, response.Code, http.StatusOK)
		require.Equal(t, response.Body.String(), fmt.Sprintf("%d", score))
	})

	t.Run("get score of user not in store", func(t *testing.T) {
		player := "Apollo"
		request := newGetScoreRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		require.Equal(t, response.Code, http.StatusNotFound)
	})
}

func TestPostgresStore_RecordWin(t *testing.T) {
	store := &PostgresStore{DB: db}
	server := NewPlayerServer(store)

	t.Run("returns accepted on POST with not present user", func(t *testing.T) {
		player := "Delta"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)

		res, err := store.DB.Query(fmt.Sprintf(`SELECT Name, Score FROM Scores WHERE name = '%s'`, player))
		require.Nil(t, err)

		defer res.Close()

		var (
			name  string
			score int
		)
		for res.Next() {
			err := res.Scan(&name, &score)
			require.Nil(t, err)
		}
		require.Equal(t, name, player)
		require.Equal(t, score, 1)

	})
	t.Run("returns accepted on POST with existing user", func(t *testing.T) {
		player := "Omega"
		scoreTest := 5
		_, err := store.DB.Exec(fmt.Sprintf(`INSERT INTO Scores (Name, Score) VALUES ('%s', %d)`, player, scoreTest))
		require.Nil(t, err)
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)

		res, err := store.DB.Query(fmt.Sprintf(`SELECT Name, Score FROM Scores WHERE name = '%s'`, player))
		require.Nil(t, err)

		defer res.Close()

		var (
			name  string
			score int
		)
		for res.Next() {
			err := res.Scan(&name, &score)
			require.Nil(t, err)
		}
		require.Equal(t, name, player)
		require.Equal(t, score, scoreTest+1)

	})
}
