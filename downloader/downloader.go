package downloader

import (
	"bytes"
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"io"
	"log"
	"multithread_downloading/config"
	"net/http"
	"os"
	"time"
)

type Downloader interface {
	DownLoad()
}

type MultiThreadDownLoader struct {
	URL        string
	ChunkSize  int
	OutputPath string
}

func (d *MultiThreadDownLoader) InitDownloader(config config.MultiThreadConfig) {
	d.URL = config.GetTarget()
	d.OutputPath = config.GetOutputPath()
	d.ChunkSize = config.ChunkSize
}

func (d *MultiThreadDownLoader) DownLoad() {
	client := &http.Client{}
	// build chunks
	chunks := BuildChunk(client, d.URL, d.ChunkSize)
	File := GetFileToSave(d.OutputPath)
	// build output channel
	WriterBlock := BuildOutputChannel()
	// download file
	go DispatchDownload(chunks, d.URL, client, WriterBlock)
	SaveInDisk(WriterBlock, d.ChunkSize, File)

}

func GetFileToSave(FilePath string) *os.File {
	file, err := os.OpenFile(FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	return file
}

func DispatchDownload(Chunks []Chunk, URL string, Client *http.Client, SaveChannel chan ChunkWriterBlock) {
	p := mpb.New(mpb.WithRefreshRate(180 * time.Millisecond))
	// download file
	for i := 0; i < len(Chunks); i++ {
		go func(i int) {
			header := http.Header{}
			header.Set("Range", "bytes="+fmt.Sprint(Chunks[i].Start)+"-"+fmt.Sprint(Chunks[i].End))
			req, err := http.NewRequest("GET", URL, nil)
			if err != nil {
				log.Fatal(err)
			}
			req.Header = header
			resp, err := Client.Do(req)
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
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			offset := Chunks[i].Start
			// write file
			var buffer bytes.Buffer
			proxyReader := bar.ProxyReader(resp.Body)
			_, err = io.Copy(&buffer, proxyReader)
			// write into channel
			SaveChannel <- ChunkWriterBlock{Buf: buffer, Offset: offset}
			if err != nil {
				log.Fatal(err)
			}
		}(i)
	}
}
