package main

import (
	"fmt"
	downloaderconfig "multithread_downloading/config/downloader"
	"multithread_downloading/downloader/multithread_downloader"
	"strings"
	"time"
)

const URL = "https://c-ssl.duitang.com/uploads/blog/202307/19/6zS5QGX8tqADG4M.jpeg"

//const URL = "https://bit.ly/1GB-testfile"

// const URL = "https://ash-speed.hetzner.com/10GB.bin"

const ChunkSize = 4

func main() {
	path := strings.Split(URL, "/")
	FilePath := path[len(path)-1]
	start := time.Now()

	downloadConfig := downloaderconfig.MultiThreadConfig{Target: URL, ChunkSize: ChunkSize, OutputPath: FilePath}
	downloaderAgent := multithread_downloader.MultiThreadDownLoader{}
	downloaderAgent.NewDownloader(downloadConfig)
	downloaderAgent.DownLoad()

	fmt.Println("Downloaded!")
	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("Time consumed: %v\n", duration)

}
