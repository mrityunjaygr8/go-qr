package main

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostgresRecordingWinsAndRetrievingThem(t *testing.T) {
	config := DB_Config{
		Host:     "localhost",
		Password: "postgres",
		Port:     15432,
		Name:     "postgres",
		Username: "postgres",
	}
	store, err := NewPostgresStore(config)
	require.Nil(t, err)
	server := PlayerServer{store}

	player := "Timon"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()

	server.ServeHTTP(response, newGetScoreRequest(player))

	require.Equal(t, response.Code, http.StatusOK)
	require.Equal(t, response.Body.String(), "3")
}
