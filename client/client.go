package client

import (
	"net/http"
	"strconv"
)

type Header struct {
	Header http.Header
}

// NewClient is the function that create a new base http client
func NewClient() *http.Client {
	client := &http.Client{}
	return client
}

// NewHeader is the function that create a new http header
func NewHeader() *Header {
	return &Header{Header: http.Header{}}
}

// HeaderAddRange is the function that add a Range value to the header
func (h *Header) HeaderAddRange(start, end int64) {
	h.Header.Add("Range", "bytes="+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10))
}

func (h *Header) GetHttpHeader() http.Header {
	return h.Header
}
