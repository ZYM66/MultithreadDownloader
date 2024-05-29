// Package multithread_downloader: This file contains the implementation of the MultiThreadDownLoader struct and its methods.
package multithread_downloader

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"io"
	c "multithread_downloading/client"
	"multithread_downloading/common"
	downloaderconfig "multithread_downloading/config/downloader"
	"multithread_downloading/storage"
	"net/http"
	"os"
	"sync"
	"time"
)

type MultiThreadDownLoader struct {
	Files        []File
	NumChunk     int
	Client       *http.Client
	HeaderConfig downloaderconfig.HeaderConfig
}

type File struct {
	Name          string
	URL           string
	OutputPath    string
	TargetFile    *os.File
	ContendLength int64
	Chunks        []Chunk
}

func NewMultiThreadDownloader(configs downloaderconfig.MultiThreadConfig) MultiThreadDownLoader {
	d := MultiThreadDownLoader{}
	d.Client = c.NewClient()
	d.NumChunk = configs.NumChunk
	d.Files = make([]File, len(configs.GetTarget()))
	d.HeaderConfig = configs.HeaderConfig

	for i, target := range configs.GetTarget() {
		f := File{}
		f.URL = target
		f.OutputPath = configs.GetOutputPath()
		f.BuildChunk(d.Client, d.NumChunk)
		f.TargetFile = storage.GetFileToSave(f.OutputPath+f.Name, f.ContendLength)
		d.Files[i] = f
	}

	return d
}

func (d *MultiThreadDownLoader) DownLoad() {
	var wg sync.WaitGroup
	closeChan := make(chan int, len(d.Files))

	for i := 0; i < len(d.Files); i++ {
		// build output channel
		OutputChannel := storage.BuildOutputChannel()

		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// download file
			d.Files[i].DispatchMultiThreadDownload(OutputChannel, d)

			// save download file into disk
			storage.SaveInDisk(OutputChannel, d.Files[i].TargetFile)

			// signal that this file is done
			closeChan <- i
		}(i)
	}

	go func() {
		wg.Wait()        // Wait for all downloads and saves to finish
		close(closeChan) // Close the channel to signal that all files are done
	}()

	// Close files as they finish
	for range closeChan {
		i := <-closeChan
		d.Files[i].TargetFile.Close()
	}
}

// DispatchMultiThreadDownload is the function that execute the MultiThreadDownload
func (f *File) DispatchMultiThreadDownload(SaveChannel chan storage.ChunkBlock, d *MultiThreadDownLoader) {
	p := mpb.New(mpb.WithRefreshRate(180 * time.Millisecond))
	// download file
	var wg sync.WaitGroup
	for i := 0; i < len(f.Chunks); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// chunk download header
			Header := c.NewHeader(d.HeaderConfig)
			Header.HeaderAddRange(f.Chunks[i].Start, f.Chunks[i].End)

			req, err := http.NewRequest("GET", f.URL, nil)
			common.Check(err)

			req.Header = Header.GetHttpHeader()
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

			offset := f.Chunks[i].Start

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
