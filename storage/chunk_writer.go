package storage

import (
	"os"
)

const WRITE_BUFFER_SIZE = 8 * 1024

type ChunkBlock struct {
	//Buf bytes.Buffer
	Buf []byte
	//BufChan chan *bytes.Buffer
	Offset int64
}

// SaveInDisk is the function that read ChunkBlock from channel then write the buffer to disk
func SaveInDisk(SaveChannel <-chan ChunkBlock, File *os.File) {
	for block := range SaveChannel {
		_, err := File.WriteAt(block.Buf, block.Offset)
		if err != nil {
			panic(err.Error())
		}
	}
}

func BuildOutputChannel() chan ChunkBlock {
	return make(chan ChunkBlock)
}
