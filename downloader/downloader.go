package downloader

import (
	"multithread_downloading/config"
)

// Downloader is the interface that define the methods that a downloader should have
type Downloader interface {
	DownLoad()
	NewDownloader(baseConfig config.DownloaderConfig)
}
