package downloader

import (
	"bytes"
	"log"
	"os"
)

type ChunkWriterBlock struct {
	Buf    bytes.Buffer
	Offset int64
}

func SaveInDisk(SaveChannel chan ChunkWriterBlock, ChunkSize int, File *os.File) {
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
