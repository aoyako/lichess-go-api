package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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

// Receive streaming response
// chan interface{} with data
// func() as cancellation
func (l *LichessAPI) streamJSON(resp *http.Response, value interface{}) (chan interface{}, func(), error) {
	reader := bufio.NewReader(resp.Body)

	finishContext, cancel := context.WithCancel(context.Background())

	receiver := make(chan interface{}, 10)

	go func() {
		defer resp.Body.Close()

		for {
			select {
			case <-finishContext.Done():
				return
			default:
				message, err := reader.ReadBytes('\n')
				if err != nil {
					close(receiver)
					return
				}
				if len(message) != 0 {
					err = json.NewDecoder(bytes.NewReader(message)).Decode(value)

					if err != nil {
						close(receiver)
						return
					}

					receiver <- reflect.ValueOf(value).Elem().Interface()
				}
			}
		}
	}()

	return receiver, cancel, nil
}

// Converts channel to chan User
func userChannelWrapper(in chan interface{}) chan User {
	middleware := make(chan User, 10)

	go func() {
		for v := range in {
			middleware <- v.(User)
		}
		close(middleware)
	}()

	return middleware
}

type reqParams struct {
	requestType string
	endpoint    string
	header      map[string]string
	query       map[string]string
	data        []byte
}
