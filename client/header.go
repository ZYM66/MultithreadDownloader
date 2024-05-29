package client

import (
	downloaderconfig "multithread_downloading/config/downloader"
	"net/http"
	"strconv"
)

type Header struct {
	Header http.Header
}

// NewHeader is the function that create a new http header
func NewHeader(config downloaderconfig.HeaderConfig) Header {
	header := Header{http.Header{}}
	if config.UA != "" {
		header.HeaderAddUA(config.UA)
	}
	return header
}

// HeaderAddRange is the function that add a Range value to the header
func (h *Header) HeaderAddRange(start, end int64) {
	h.Header.Add("Range", "bytes="+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10))
}

func (h *Header) HeaderAddUA(ua string) {
	h.Header.Add("User-Agent", ua)
}

func (h *Header) GetHttpHeader() http.Header {
	return h.Header
}
