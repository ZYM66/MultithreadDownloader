package storage

import (
	"fmt"
	"log"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetFileToSave is the function that create a file in disk to save the bytes got from network
func GetFileToSave(FilePath string) *os.File {
	if ok, _ := PathExists(FilePath); !ok {
		file, err := os.OpenFile(FilePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err.Error())
		}
		return file
	} else {
		stat, _ := os.Stat(FilePath)
		fmt.Println(stat.Size())
	}

	return nil
}
