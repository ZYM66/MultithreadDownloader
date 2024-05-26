package multithread_downloader

import (
	"log"
	"net/http"
)

type Chunk struct {
	Start int64
	End   int64
}

func BuildChunk(client *http.Client, url string, numChunks int) []Chunk {
	// Make a GET request to the URL
	res, err := client.Get(url)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()
	// Get the content length
	length := res.ContentLength
	// Calculate chunk size
	chunkSize := length / int64(numChunks)
	// Initialize chunks slice
	chunks := make([]Chunk, numChunks)
	// Build each chunk
	for i := 0; i < numChunks; i++ {
		// Calculate start and end positions for the chunk
		start := int64(i) * chunkSize
		end := start + chunkSize - 1 // Subtract 1 to avoid overlap
		// For the last chunk, set the end to the length of the content
		if i == numChunks-1 {
			end = length - 1
		}
		// Create a Chunk instance and add it to the chunks slice
		chunks[i] = Chunk{Start: start, End: end}
	}

	return chunks
}
