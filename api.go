package main

import (
	"bytes"
	"fmt"
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
	client      *http.Client
	endpoint    *serviceEndpoint
}

// NewLichessAPI creates LichessAPI struct from config
func NewLichessAPI(cfg Config) *LichessAPI {
	apiController := &LichessAPI{
		bearerToken: fmt.Sprintf("Bearer %s", cfg.Token),
		client:      cfg.Client,
		endpoint:    newServiceEndpoint(),
	}

	if apiController.client == nil {
		apiController.client = &http.Client{}
	}

	return apiController
}

// request is called when requesting data from lichess.org
func (l *LichessAPI) request(par *reqParams) (*http.Response, error) {
	req, err := http.NewRequest(par.requestType, par.endpoint, bytes.NewBuffer(par.data))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", l.bearerToken)

	if par.header != nil {
		for k, v := range par.header {
			req.Header.Add(k, v)
		}
	}

	if par.query != nil {
		q := req.URL.Query()
		for k, v := range par.query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type reqParams struct {
	requestType string
	endpoint    string
	header      map[string]string
	query       map[string]string
	data        []byte
}
