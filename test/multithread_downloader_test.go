package testpackage

import (
	downloaderconfig "multithread_downloading/config/downloader"
	"multithread_downloading/downloader/multithread_downloader"
	"testing"
)

func TestNewDownloader(t *testing.T) {
	downloadConfig := downloaderconfig.MultiThreadConfig{Target: "https://example.com", ChunkSize: 4, OutputPath: "output.txt"}
	downloaderAgent := multithread_downloader.MultiThreadDownLoader{}
	downloaderAgent.NewDownloader(downloadConfig)
}

func TestDownLoad(t *testing.T) {
	// This test would require setting up a local server to serve a test file
	// and then checking that the file was correctly downloaded.
}

func TestDownLoadInvalidURL(t *testing.T) {
	// This test would check that the DownLoad method correctly handles an invalid URL
}

func TestDownLoadLargeChunkSize(t *testing.T) {
	// This test would check that the DownLoad method correctly handles a large chunk size
}

func TestDownLoadNonExistentOutputPath(t *testing.T) {
	// This test would check that the DownLoad method correctly handles a non-existent output path
}
