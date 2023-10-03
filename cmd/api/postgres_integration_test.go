package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostgresRecordingWinsAndRetrievingThem(t *testing.T) {
	config := DbConfig{
		Host:     "localhost",
		Password: "postgres",
		Port:     "15432",
		Name:     "postgres",
		Username: "postgres",
	}
	store, err := NewPostgresStore(config)
	require.Nil(t, err)
	server := NewPlayerServer(store)

	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()

		server.ServeHTTP(response, newGetScoreRequest(player))

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()

		server.ServeHTTP(response, newLeagueRequest())
		require.Equal(t, http.StatusOK, response.Code)

		var got []Player
		err := json.NewDecoder(response.Body).Decode(&got)
		require.Nil(t, err)

		want := []Player{
			{"Pepper", 3},
		}

		require.Equal(t, want, got)
	})
}
