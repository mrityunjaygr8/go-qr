package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "db")
	assert.Nilf(t, err, "could not create temp file %v", err)

	tmpFile.Write([]byte(initialData))

	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return tmpFile, removeFile
}

func TestFileSystemPlayerStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name": "Cleo", "Wins": 10},
				{"Name": "Chris", "Wins": 20}
			]`)

		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assert.Nil(t, err)

		got := store.GetLeague()
		want := League{{"Cleo", 10}, {"Chris", 20}}

		assert.Equal(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name": "Cleo", "Wins": 10},
				{"Name": "Chris", "Wins": 20}
			]`)

		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assert.Nil(t, err)

		got := store.GetPlayerScore("Chris")
		want := 20

		assert.Equal(t, want, got)
	})

	t.Run("record score for existing user", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name": "Cleo", "Wins": 10},
				{"Name": "Chris", "Wins": 20}
			]`)

		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assert.Nil(t, err)

		store.RecordWin("Chris")
		got := store.GetPlayerScore("Chris")
		want := 21

		assert.Equal(t, want, got)
	})

	t.Run("record score for new user", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name": "Cleo", "Wins": 10},
				{"Name": "Chris", "Wins": 20}
			]`)

		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assert.Nil(t, err)

		store.RecordWin("Kurt")
		got := store.GetPlayerScore("Kurt")
		want := 1

		assert.Equal(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)
		assert.Nil(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 20}
		]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assert.Nil(t, err)

		got := store.GetLeague()

		want := League{
			{"Chris", 20},
			{"Cleo", 10},
		}

		assert.Equal(t, want, got)
	})
}
