package main

import (
	"fmt"
	"net/http"
)

// GetFollowing returns users which are followed by specified user
func (l *LichessAPI) GetFollowing(id string) (chan User, func(), error) {
	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    fmt.Sprintf(l.endpoint.relationFollowing, id),
	}

	resp, err := l.request(params)
	if err != nil {
		return nil, nil, err
	}

	uchan, cancel, err := l.streamJSON(resp, &User{})

	return userChannelWrapper(uchan), cancel, err
}

// GetFollowers returns users who follow specified user
func (l *LichessAPI) GetFollowers(id string) (chan User, func(), error) {
	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    fmt.Sprintf(l.endpoint.relationFollowers, id),
	}

	resp, err := l.request(params)
	if err != nil {
		return nil, nil, err
	}

	uchan, cancel, err := l.streamJSON(resp, &User{})

	return userChannelWrapper(uchan), cancel, err
}
