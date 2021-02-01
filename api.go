package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// Config stores account configuration
type Config struct {
	Token  string
	Client *http.Client
}

// LichessAPI represents struct that handles api calls
type LichessAPI struct {
	bearerToken string
	client      http.Client
	endpoint    *serviceEndpoint
}

// NewLichessAPI creates LichessAPI struct from config
func NewLichessAPI(cfg Config) *LichessAPI {
	return &LichessAPI{
		bearerToken: fmt.Sprintf("Bearer %s", cfg.Token),
		client:      *cfg.Client,
		endpoint:    newServiceEndpoint(),
	}
}

// request is called when requesting data from lichess.org
func (l *LichessAPI) request(rtype, endpoint, data string) *http.Response {
	req, err := http.NewRequest(rtype, endpoint, bytes.NewBuffer([]byte(data)))
	req.Header.Add("Authorization", l.bearerToken)

	resp, err := l.client.Do(req)
	if err != nil {
		log.Printf("Lichess request error: %s\n", err)
	}

	return resp
}
