package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetFollowers(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		ID string
	}
	type params struct {
		requestType string
		requestBody string
		response    string
		query       string
	}
	tests := []struct {
		name    string
		args    args
		params  params
		want    []User
		wantErr error
	}{
		{
			name: "Get followers",
			args: args{
				ID: "user",
			},
			params: params{
				requestType: http.MethodGet,
				response: "{\"id\": \"bahadirozen\", \"username\": \"BahadirOzen\"}\n" +
					"{\"id\": \"penguingim1\", \"username\": \"penguingim1\"}\n",
			},
			want: []User{
				{
					ID:       "bahadirozen",
					Username: "BahadirOzen",
				},
				{
					ID:       "penguingim1",
					Username: "penguingim1",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(tt.params.requestType, req.Method)

				body, _ := ioutil.ReadAll(req.Body)
				assert.Equal(tt.params.requestBody, string(body))

				rw.Write([]byte(tt.params.response))
			}))

			lapi := NewLichessAPI(Config{
				Token:  "",
				Client: server.Client(),
			})
			lapi.endpoint.relationFollowers = server.URL + "/%s"

			uchan, _, err := lapi.GetFollowers(tt.args.ID)

			var users []User

			for u := range uchan {
				users = append(users, u)
			}

			assert.Equal(tt.want, users)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_GetFollowing(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		ID string
	}
	type params struct {
		requestType string
		requestBody string
		response    string
		query       string
	}
	tests := []struct {
		name    string
		args    args
		params  params
		want    []User
		wantErr error
	}{
		{
			name: "Get followers",
			args: args{
				ID: "user",
			},
			params: params{
				requestType: http.MethodGet,
				response: "{\"id\": \"bahadirozen\", \"username\": \"BahadirOzen\"}\n" +
					"{\"id\": \"penguingim1\", \"username\": \"penguingim1\"}\n",
			},
			want: []User{
				{
					ID:       "bahadirozen",
					Username: "BahadirOzen",
				},
				{
					ID:       "penguingim1",
					Username: "penguingim1",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(tt.params.requestType, req.Method)

				body, _ := ioutil.ReadAll(req.Body)
				assert.Equal(tt.params.requestBody, string(body))

				rw.Write([]byte(tt.params.response))
			}))

			lapi := NewLichessAPI(Config{
				Token:  "",
				Client: server.Client(),
			})
			lapi.endpoint.relationFollowing = server.URL + "/%s"

			uchan, _, err := lapi.GetFollowing(tt.args.ID)

			var users []User

			for u := range uchan {
				users = append(users, u)
			}

			assert.Equal(tt.want, users)
			assert.Equal(tt.wantErr, err)
		})
	}
}
