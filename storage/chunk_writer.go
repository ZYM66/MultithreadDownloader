package storage

import (
	"bytes"
	"log"
	"os"
)

type ChunkWriterBlock struct {
	Buf    bytes.Buffer
	Offset int64
}

// SaveInDisk is the function that read ChunkWriterBlock from channel then write the buffer to disk
func SaveInDisk(SaveChannel <-chan ChunkWriterBlock, ChunkSize int, File *os.File) {
	for i := 0; i < ChunkSize; i++ {
		block, _ := <-SaveChannel

		_, err := File.WriteAt(block.Buf.Bytes(), block.Offset)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func BuildOutputChannel() chan ChunkWriterBlock {
	return make(chan ChunkWriterBlock)
}
