package client

import (
	"net/http"
)

// NewClient is the function that create a new base http client
func NewClient() *http.Client {
	client := &http.Client{}
	return client
}
