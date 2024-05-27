// Package multithread_downloader: This file contains the implementation of the MultiThreadDownLoader struct and its methods.
package multithread_downloader

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"io"
	"log"
	c "multithread_downloading/client"
	"multithread_downloading/common"
	"multithread_downloading/config"
	downloaderconfig "multithread_downloading/config/downloader"
	"multithread_downloading/storage"
	"net/http"
	"sync"
	"time"
)

type MultiThreadDownLoader struct {
	URL        string
	ChunkSize  int
	OutputPath string
}

func NewMultiThreadDownloader(configs config.DownloaderConfig) MultiThreadDownLoader {
	d := MultiThreadDownLoader{}
	if v, ok := configs.(downloaderconfig.MultiThreadConfig); ok {
		d.ChunkSize = v.ChunkSize
	} else {
		log.Fatal("Invalid config")
	}
	d.URL = configs.GetTarget()
	d.OutputPath = configs.GetOutputPath()
	return d
}

func (d *MultiThreadDownLoader) DownLoad() {
	// build client
	client := c.NewClient()
	// build chunks
	chunks := BuildChunk(client, d.URL, d.ChunkSize)
	// create file in disk
	File := storage.GetFileToSave(d.OutputPath)
	// build output channel
	OutputChannel := storage.BuildOutputChannel()
	// download file
	go DispatchMultiThreadDownload(chunks, d.URL, client, OutputChannel)
	// save download file into disk
	storage.SaveInDisk(OutputChannel, File)

	defer File.Close()

}

// DispatchMultiThreadDownload is the function that execute the MultiThreadDownload
func DispatchMultiThreadDownload(Chunks []Chunk, URL string, Client *http.Client, SaveChannel chan storage.ChunkWriterBlock) {
	p := mpb.New(mpb.WithRefreshRate(180 * time.Millisecond))
	// download file
	var wg sync.WaitGroup
	for i := 0; i < len(Chunks); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			header := c.NewHeader()
			header.HeaderAddRange(Chunks[i].Start, Chunks[i].End)
			req, err := http.NewRequest("GET", URL, nil)
			common.Check(err)

			req.Header = header.GetHttpHeader()
			resp, err := Client.Do(req)
			common.Check(err)

			// Check the status code and Content-Range header
			for {
				if resp.StatusCode != http.StatusPartialContent {
					resp, err = Client.Do(req)
				} else {
					break
				}
			}

			//add progressbar
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
			proxyReader := bar.ProxyReader(resp.Body)
			buf := make([]byte, 32*1024) // 32KB buffer
			for {
				n, err := proxyReader.Read(buf)
				if err != nil {
					if err == io.EOF {
						break
					} else {
						common.Check(err)
					}
				}
				// 深拷贝 buf[:n] 到一个新的切片中，确保传递给 SaveChannel 的数据独立
				SaveChannel <- storage.ChunkWriterBlock{Buf: append([]byte{}, buf[:n]...), Offset: offset}
				offset += int64(n)
			}
		}(i)

	}
	// close the SaveChannel
	go func() {
		wg.Wait()
		close(SaveChannel)
	}()
}
