package storage

import (
	"log"
	"os"
)

const WRITE_BUFFER_SIZE = 8 * 1024

type ChunkWriterBlock struct {
	//Buf bytes.Buffer

	Buf []byte
	//BufChan chan *bytes.Buffer
	Offset int64
}

// SaveInDisk is the function that read ChunkWriterBlock from channel then write the buffer to disk
func SaveInDisk(SaveChannel <-chan ChunkWriterBlock, File *os.File) {
	for block := range SaveChannel {
		//fmt.Println(block.Offset, len(block.Buf))
		_, err := File.WriteAt(block.Buf, block.Offset)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func BuildOutputChannel() chan ChunkWriterBlock {
	return make(chan ChunkWriterBlock)
}
