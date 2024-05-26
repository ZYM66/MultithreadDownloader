// Package multithread_downloader: This file contains the implementation of the MultiThreadDownLoader struct and its methods.
package multithread_downloader

import (
	"bytes"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"io"
	"log"
	c "multithread_downloading/client"
	"multithread_downloading/common"
	"multithread_downloading/config"
	downloaderconfig "multithread_downloading/config/downloader"
	"multithread_downloading/downloader"
	"multithread_downloading/storage"
	"net/http"
	"time"
)

type MultiThreadDownLoader struct {
	URL        string
	ChunkSize  int
	OutputPath string
}

func (d *MultiThreadDownLoader) NewDownloader(configs config.DownloaderConfig) {
	if v, ok := configs.(downloaderconfig.MultiThreadConfig); ok {
		d.ChunkSize = v.ChunkSize
	} else {
		log.Fatal("Invalid config")
	}
	d.URL = configs.GetTarget()
	d.OutputPath = configs.GetOutputPath()

}

func (d *MultiThreadDownLoader) DownLoad() {
	// build client
	client := c.NewClient()
	//client := &http.Client{}
	// build chunks
	chunks := downloader.BuildChunk(client, d.URL, d.ChunkSize)
	// create file in disk
	File := storage.GetFileToSave(d.OutputPath)
	// build output channel
	WriterBlock := storage.BuildOutputChannel()
	// download file
	go DispatchMultiThreadDownload(chunks, d.URL, client, WriterBlock)
	// save download file into disk
	storage.SaveInDisk(WriterBlock, d.ChunkSize, File)

	defer File.Close()

}

// DispatchMultiThreadDownload is the function that execute the MultiThreadDownload
func DispatchMultiThreadDownload(Chunks []downloader.Chunk, URL string, Client *http.Client, SaveChannel chan storage.ChunkWriterBlock) {
	p := mpb.New(mpb.WithRefreshRate(180 * time.Millisecond))
	// download file
	for i := 0; i < len(Chunks); i++ {
		go func(i int) {
			header := c.NewHeader()
			header.HeaderAddRange(Chunks[i].Start, Chunks[i].End)
			req, err := http.NewRequest("GET", URL, nil)
			common.Check(err)

			req.Header = header.GetHttpHeader()
			resp, err := Client.Do(req)
			common.Check(err)

			//add progressbar
			// Check the status code and Content-Range header
			for {
				if resp.StatusCode != http.StatusPartialContent {
					resp, err = Client.Do(req)
				} else {
					break
				}
			}

			bar := p.New(resp.ContentLength,
				mpb.BarStyle().Rbound("|"),
				mpb.PrependDecorators(
					decor.Counters(decor.SizeB1024(0), "% .3f / % .3f"),
				),
				mpb.AppendDecorators(
					decor.EwmaETA(decor.ET_STYLE_GO, 30),
					decor.Name(" ] "),
					decor.EwmaSpeed(decor.SizeB1024(0), "% .2f", 60),
				),
			)
			common.Check(err)
			defer resp.Body.Close()

			offset := Chunks[i].Start

			// write file
			var buffer bytes.Buffer
			proxyReader := bar.ProxyReader(resp.Body)
			_, err = io.Copy(&buffer, proxyReader)
			common.Check(err)

			// write into channel
			SaveChannel <- storage.ChunkWriterBlock{Buf: buffer, Offset: offset}
		}(i)
	}
}
