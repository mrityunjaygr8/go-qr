package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func Test_tape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &tape{file: file}
	tape.Write([]byte("abc"))

	file.Seek(0, 0)
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	assert.Equal(t, want, got)
}
