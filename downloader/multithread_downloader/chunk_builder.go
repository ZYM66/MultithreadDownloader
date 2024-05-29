package multithread_downloader

import (
	"log"
	"net/http"
	"strings"
)

type Chunk struct {
	Start int64
	End   int64
}

func (f *File) BuildChunk(client *http.Client, num_chunk int) {
	// Make a GET request to the URL
	res, err := client.Get(f.URL)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()
	// Get the content length
	length := res.ContentLength
	// get file name
	fileName := res.Header.Get("Content-Disposition")
	if fileName == "" {
		fileName = strings.Split(res.Request.URL.Path, "/")[len(strings.Split(res.Request.URL.Path, "/"))-1]
	}
	f.Name = fileName

	// Set the content length in the downloader instance
	f.ContendLength = length
	// Calculate chunk size
	chunkSize := length / int64(num_chunk)
	// Initialize chunks slice
	chunks := make([]Chunk, num_chunk)
	// Build each chunk
	for i := 0; i < num_chunk; i++ {
		// Calculate start and end positions for the chunk
		start := int64(i) * chunkSize
		end := start + chunkSize - 1 // Subtract 1 to avoid overlap
		// For the last chunk, set the end to the length of the content
		if i == num_chunk-1 {
			end = length - 1
		}
		// Create a Chunk instance and add it to the chunks slice
		chunks[i] = Chunk{Start: start, End: end}
	}

	f.Chunks = chunks
}
