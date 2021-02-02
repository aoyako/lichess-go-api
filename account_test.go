package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetMyProfile(t *testing.T) {
	assert := assert.New(t)

	type params struct {
		path        string
		requestType string
		requestBody string
		response    string
	}
	tests := []struct {
		name    string
		params  params
		want    User
		wantErr error
	}{
		{
			name: "Get user data",
			params: params{
				path:        "https://lichess.org/api/account",
				requestType: http.MethodGet,
				requestBody: "",
				response: `{
					"id": "georges",
					"username": "Georges",
					"online": true,
					"perfs": {
					  "blitz": {
						"games": 2945,
						"rating": 1609,
						"rd": 60,
						"prog": -22,
						"prov": true
					  },
					  "puzzle": {
						"games": 2945,
						"rating": 1609,
						"rd": 60,
						"prog": -22,
						"prov": true
					  }
					},
					"createdAt": 1290415680000,
					"disabled": false,
					"tosViolation": false,
					"booster": false,
					"profile": {
					  "country": "EC",
					  "location": "string",
					  "bio": "Free bugs!",
					  "firstName": "Thibault",
					  "lastName": "Duplessis",
					  "fideRating": 1500,
					  "uscfRating": 1500,
					  "ecfRating": 1500,
					  "links": "github.com/ornicar\r\ntwitter.com/ornicar"
					},
					"seenAt": 1522636452014,
					"patron": true,
					"playTime": {
					  "total": 3296897,
					  "tv": 12134
					},
					"language": "en-GB",
					"title": "NM",
					"url": "https://lichess.org/@/georges",
					"playing": "https://lichess.org/yqfLYJ5E/black",
					"nbFollowing": 299,
					"nbFollowers": 2735,
					"completionRate": 97,
					"count": {
					  "all": 9265,
					  "rated": 7157,
					  "ai": 531,
					  "draw": 340,
					  "drawH": 331,
					  "loss": 4480,
					  "lossH": 4207,
					  "win": 4440,
					  "winH": 4378,
					  "bookmark": 71,
					  "playing": 6,
					  "import": 66,
					  "me": 0
					},
					"streaming": false,
					"followable": true,
					"following": false,
					"blocking": false,
					"followsYou": false
				  }`,
			},
			want: User{
				ID:       "georges",
				Username: "Georges",
				Online:   true,
				Performances: map[string]Performance{
					"blitz": {
						Games:           2945,
						Rating:          1609,
						RatingDeviation: 60,
						Prog:            -22,
						Prov:            true,
					},
					"puzzle": {
						Games:           2945,
						Rating:          1609,
						RatingDeviation: 60,
						Prog:            -22,
						Prov:            true,
					},
				},
				CreatedAt:    1290415680000,
				Disabled:     false,
				TOSViolation: false,
				Booster:      false,
				Profile: Profile{
					Country:    "EC",
					Location:   "string",
					Bio:        "Free bugs!",
					FirstName:  "Thibault",
					LastName:   "Duplessis",
					FideRating: 1500,
					UscfRating: 1500,
					EcfRating:  1500,
					Links:      "github.com/ornicar\r\ntwitter.com/ornicar",
				},
				SeenAt: 1522636452014,
				Patron: true,
				PlayTime: PlayTime{
					Total: 3296897,
					TV:    12134,
				},
				Language:        "en-GB",
				Title:           "NM",
				URL:             "https://lichess.org/@/georges",
				Playing:         "https://lichess.org/yqfLYJ5E/black",
				NumberFollowing: 299,
				NumberFollowers: 2735,
				CompletionRate:  97,
				Count: StatsCount{
					All:      9265,
					Rated:    7157,
					AI:       531,
					Draw:     340,
					DrawH:    331,
					Loss:     4480,
					LossH:    4207,
					Win:      4440,
					WinH:     4378,
					Bookmark: 71,
					Playing:  6,
					Import:   66,
					Me:       0,
				},
				Streaming:  false,
				Followable: true,
				Following:  false,
				Blocking:   false,
				FollowsYou: false,
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
			lapi.endpoint.accountProfile = server.URL

			user, err := lapi.GetMyProfile()

			assert.Equal(tt.want, *user)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_GetMyEmail(t *testing.T) {
	assert := assert.New(t)

	type params struct {
		path        string
		requestType string
		requestBody string
		response    string
	}
	tests := []struct {
		name    string
		params  params
		want    string
		wantErr error
	}{
		{
			name: "Get user email",
			params: params{
				requestType: http.MethodGet,
				requestBody: "",
				response: `{
					"email": "example@nomail.org"
				}`,
			},
			want:    "example@nomail.org",
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
			lapi.endpoint.accountEmail = server.URL

			mail, err := lapi.GetMyEmail()

			assert.Equal(tt.want, mail)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_GetMyKidStatus(t *testing.T) {
	assert := assert.New(t)

	type params struct {
		path        string
		requestType string
		requestBody string
		response    string
	}
	tests := []struct {
		name    string
		params  params
		want    bool
		wantErr error
	}{
		{
			name: "Get user kid status",
			params: params{
				requestType: http.MethodGet,
				requestBody: "",
				response: `{
					"kid": true
				}`,
			},
			want:    true,
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
			lapi.endpoint.accountKidModeStatus = server.URL

			kid, err := lapi.GetMyKidModeStatus()

			assert.Equal(tt.want, kid)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_SetMyKidStatus(t *testing.T) {
	assert := assert.New(t)

	type params struct {
		path        string
		requestType string
		requestBody string
		response    string
	}
	tests := []struct {
		name    string
		args    bool
		params  params
		want    bool
		wantErr error
	}{
		{
			name: "Set user kid status",
			args: true,
			params: params{
				requestType: http.MethodPost,
				requestBody: "v=true",
				response: `{
					"ok": true
				}`,
			},
			want:    true,
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
			lapi.endpoint.accountKidModeStatus = server.URL

			kid, err := lapi.SetMyKidModeStatus(tt.args)

			assert.Equal(tt.want, kid)
			assert.Equal(tt.wantErr, err)
		})
	}
}
