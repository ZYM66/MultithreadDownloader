package storage

import (
	"log"
	"os"
)

// GetFileToSave is the function that create a file in disk to save the bytes got from network
func GetFileToSave(FilePath string) *os.File {
	file, err := os.OpenFile(FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	return file
}
