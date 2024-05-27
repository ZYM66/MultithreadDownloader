package storage

import (
	"fmt"
	"multithread_downloading/common"
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
func GetFileToSave(FilePath string, ContendLength int64) *os.File {
	if ok, _ := PathExists(FilePath); !ok {
		file, err := os.OpenFile(FilePath, os.O_WRONLY|os.O_CREATE, 0644)
		common.Check(err)
		// set the file size
		err = file.Truncate(ContendLength)
		fmt.Println(ContendLength)
		common.Check(err)
		return file
	} else {
		stat, _ := os.Stat(FilePath)
		fmt.Println(stat.Size())
	}

	return nil
}
