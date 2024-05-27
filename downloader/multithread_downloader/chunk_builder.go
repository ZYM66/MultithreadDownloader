package multithread_downloader

import (
	"log"
)

type Chunk struct {
	Start int64
	End   int64
}

func (d *MultiThreadDownLoader) BuildChunk() {
	// Make a GET request to the URL
	res, err := d.Client.Get(d.URL)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()
	// Get the content length
	length := res.ContentLength
	// Set the content length in the downloader instance
	d.ContendLength = length
	// Calculate chunk size
	chunkSize := length / int64(d.NumChunk)
	// Initialize chunks slice
	chunks := make([]Chunk, d.NumChunk)
	// Build each chunk
	for i := 0; i < d.NumChunk; i++ {
		// Calculate start and end positions for the chunk
		start := int64(i) * chunkSize
		end := start + chunkSize - 1 // Subtract 1 to avoid overlap
		// For the last chunk, set the end to the length of the content
		if i == d.NumChunk-1 {
			end = length - 1
		}
		// Create a Chunk instance and add it to the chunks slice
		chunks[i] = Chunk{Start: start, End: end}
	}

	d.Chunks = chunks
}
