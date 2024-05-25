package main

import (
	"fmt"
	"multithread_downloading/config"
	"multithread_downloading/downloader"
	"strings"
	"time"
)

const URL = "https://c-ssl.duitang.com/uploads/blog/202307/19/6zS5QGX8tqADG4M.jpeg"

// const url = "https://bit.ly/1GB-testfile"
// const URL = "https://ash-speed.hetzner.com/10GB.bin"
const ChunkSize = 8

func main() {
	path := strings.Split(URL, "/")
	FilePath := path[len(path)-1]
	start := time.Now()

	downloaderAgent := &downloader.MultiThreadDownLoader{}
	downloaderAgent.InitDownloader(config.MultiThreadConfig{
		BaseConfig: config.BaseConfig{Target: URL, OutputPath: FilePath},
		ChunkSize:  ChunkSize,
	})
	downloaderAgent.DownLoad()

	fmt.Println("Downloaded!")
	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("Time consumed: %v\n", duration)

}
