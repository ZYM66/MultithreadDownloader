package main

import (
	"flag"
	"fmt"
	"log"
	"multithread_downloading/common"
	downloaderconfig "multithread_downloading/config/downloader"
	"multithread_downloading/downloader/multithread_downloader"
	"multithread_downloading/statistic"
	"strings"
	"time"
)

// const URL = "https://c-ssl.duitang.com/uploads/blog/202307/19/6zS5QGX8tqADG4M.jpeg"
// const URL = "https://bit.ly/1GB-testfile"

// const URL = "https://ash-speed.hetzner.com/10GB.bin"
var (
	URLs      []string
	PATH      string
	ChunkSize int
)

func init() {
	flag.Var((*common.StringSlice)(&URLs), "url", "file urls separated by comma")
	flag.StringVar(&PATH, "path", "./", "file save path")
	flag.IntVar(&ChunkSize, "num_thread", 4, "number of download thread")
	flag.Parse()
}

// todo: 下载完成后sha256校验
// todo: 断点续传

func main() {

	defer statistic.TimeCost(time.Now())
	if len(URLs) == 0 {
		log.Fatal("Invalid URL")
	} else {
		for i := range URLs {
			URLs[i] = strings.TrimSpace(URLs[i])
		}
	}

	downloadConfig := downloaderconfig.MultiThreadConfig{Target: URLs, NumChunk: ChunkSize, OutputPath: PATH, HeaderConfig: downloaderconfig.HeaderConfig{UA: "netdisk;PC"}}
	downloaderAgent := multithread_downloader.NewMultiThreadDownloader(downloadConfig)
	downloaderAgent.DownLoad()

	fmt.Println("Downloaded!")
}
