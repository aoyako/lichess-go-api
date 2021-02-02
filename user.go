package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GetUserStatus returns user's status from their ids
func (l *LichessAPI) GetUserStatus(ids ...string) ([]User, error) {
	var users []User

	if len(ids) == 0 {
		return users, nil
	}
	if len(ids) > 50 {
		return users, errors.New("Too many requested profiles, max is 50")
	}

	var sb strings.Builder

	for pos := range ids[:len(ids)-1] {
		sb.WriteString(ids[pos])
		sb.WriteString(", ")
	}
	sb.WriteString(ids[len(ids)-1])

	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    l.endpoint.userStatus,
		query: map[string]string{
			"ids": sb.String(),
		},
	}

	resp, err := l.request(params)
	if err != nil {
		return users, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return users, err
	}

	return users, nil
}

// GetAllTop returns map with category and top 10 players
func (l *LichessAPI) GetAllTop() (map[string][]User, error) {
	var top map[string][]User

	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    l.endpoint.topAllPlayers,
		header: map[string]string{
			"Accept": "application/vnd.lichess.v3+json",
		},
	}

	resp, err := l.request(params)
	if err != nil {
		return top, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&top)
	if err != nil {
		return top, err
	}

	return top, nil
}

// GetTop returns top N users in specified game category
func (l *LichessAPI) GetTop(category string, number int) ([]User, error) {
	type users struct {
		Users []User `json:"users"`
	}

	var top users

	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    fmt.Sprintf(l.endpoint.topPlayers, number, category),
		header: map[string]string{
			"Accept": "application/vnd.lichess.v3+json",
		},
	}

	resp, err := l.request(params)
	if err != nil {
		return top.Users, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&top)
	if err != nil {
		return top.Users, err
	}

	return top.Users, nil
}

// GetUser returns public information about user with specified username
func (l *LichessAPI) GetUser(username string) (*User, error) {
	var user User

	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    fmt.Sprintf(l.endpoint.userProfile, username),
	}

	resp, err := l.request(params)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
