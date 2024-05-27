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
	"os"
	"sync"
	"time"
)

type MultiThreadDownLoader struct {
	URL           string
	NumChunk      int
	OutputPath    string
	TargetFile    *os.File
	ContendLength int64
	Chunks        []Chunk
	Client        *http.Client
}

func NewMultiThreadDownloader(configs config.DownloaderConfig) MultiThreadDownLoader {
	d := MultiThreadDownLoader{}
	if v, ok := configs.(downloaderconfig.MultiThreadConfig); ok {
		d.NumChunk = v.NumChunk
	} else {
		log.Fatal("Invalid config")
	}
	d.URL = configs.GetTarget()
	d.OutputPath = configs.GetOutputPath()
	// build client
	d.Client = c.NewClient()
	d.BuildChunk()
	d.TargetFile = storage.GetFileToSave(d.OutputPath, d.ContendLength)
	return d
}

func (d *MultiThreadDownLoader) DownLoad() {
	// build output channel
	OutputChannel := storage.BuildOutputChannel()
	// download file
	go d.DispatchMultiThreadDownload(OutputChannel)
	// save download file into disk
	storage.SaveInDisk(OutputChannel, d.TargetFile)

	defer d.TargetFile.Close()

}

// DispatchMultiThreadDownload is the function that execute the MultiThreadDownload
func (d *MultiThreadDownLoader) DispatchMultiThreadDownload(SaveChannel chan storage.ChunkBlock) {
	p := mpb.New(mpb.WithRefreshRate(180 * time.Millisecond))
	// download file
	var wg sync.WaitGroup
	for i := 0; i < len(d.Chunks); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			header := c.NewHeader()
			header.HeaderAddRange(d.Chunks[i].Start, d.Chunks[i].End)
			req, err := http.NewRequest("GET", d.URL, nil)
			common.Check(err)

			req.Header = header.GetHttpHeader()
			resp, err := d.Client.Do(req)
			common.Check(err)

			// Check the status code and Content-Range header, if not 206 Partial Content, re-send the request
			for {
				if resp.StatusCode != http.StatusPartialContent {
					resp, err = d.Client.Do(req)
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

			offset := d.Chunks[i].Start

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
				SaveChannel <- storage.ChunkBlock{Buf: append([]byte{}, buf[:n]...), Offset: offset}
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
