package main

import (
	"flag"
	"fmt"
	"log"
	downloaderconfig "multithread_downloading/config/downloader"
	"multithread_downloading/downloader/multithread_downloader"
	"multithread_downloading/statistic"
	"strings"
	"time"
)

// const URL = "https://c-ssl.duitang.com/uploads/blog/202307/19/6zS5QGX8tqADG4M.jpeg"
//const URL = "https://bit.ly/1GB-testfile"

// const URL = "https://ash-speed.hetzner.com/10GB.bin"
var (
	URL       string
	PATH      string
	ChunkSize int
)

func init() {
	flag.StringVar(&URL, "url", "", "file url")
	flag.StringVar(&PATH, "path", "", "file save path")
	flag.IntVar(&ChunkSize, "num_thread", 4, "number of download thread")
	flag.Parse()
}

func main() {

	defer statistic.TimeCost(time.Now())
	if URL == "" {
		log.Fatal("Invalid URL")
	}
	if PATH == "" {
		p := strings.Split(URL, "/")
		PATH = p[len(p)-1]
	}
	downloadConfig := downloaderconfig.MultiThreadConfig{Target: URL, ChunkSize: ChunkSize, OutputPath: PATH}
	downloaderAgent := multithread_downloader.NewMultiThreadDownloader(downloadConfig)
	downloaderAgent.DownLoad()

	fmt.Println("Downloaded!")
}
