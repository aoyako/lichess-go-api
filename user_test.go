package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetUserStatus(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		IDs []string
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
			name: "Get user status",
			args: args{
				IDs: []string{"user1", "user2"},
			},
			params: params{
				query:       "user1,user2",
				requestType: http.MethodGet,
				requestBody: "",
				response: `[{
					"id": "user1",
					"username": "user1",
					"online": true
				  },
				  {
					"id": "user2",
					"username": "user2",
					"online": true
				  }]`,
			},
			want: []User{
				{
					ID:       "user1",
					Username: "user1",
					Online:   true,
				},
				{
					ID:       "user2",
					Username: "user2",
					Online:   true,
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

				keys, ok := req.URL.Query()["ids"]
				assert.True(ok)
				assert.Equal(tt.params.query, keys[0])

				rw.Write([]byte(tt.params.response))
			}))

			lapi := NewLichessAPI(Config{
				Token:  "",
				Client: server.Client(),
			})
			lapi.endpoint.userStatus = server.URL

			users, err := lapi.GetUserStatus(tt.args.IDs...)

			assert.Equal(tt.want, users)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_GetAllTop(t *testing.T) {
	assert := assert.New(t)

	type params struct {
		requestType string
		requestBody string
		response    string
		query       string
	}
	tests := []struct {
		name    string
		params  params
		want    map[string][]User
		wantErr error
	}{
		{
			name: "Get all top",
			params: params{
				requestType: http.MethodGet,
				requestBody: "",
				response: `{
					"bullet": [
						{
							"id": "bahadirozen",
							"username": "BahadirOzen"
						},
						{
							"id": "penguingim1",
							"username": "penguingim1"
						}
					],
					"blitz": [],
					"rapid": [],
					"classical": [],
					"ultraBullet": [],
					"chess960": [],
					"crazyhouse": [],
					"antichess": [],
					"atomic": [],
					"horde": [],
					"kingOfTheHill": [],
					"racingKings": [],
					"threeCheck": []
				
				}`,
			},
			want: map[string][]User{
				"bullet": {
					{
						ID:       "bahadirozen",
						Username: "BahadirOzen",
					},
					{
						ID:       "penguingim1",
						Username: "penguingim1",
					},
				},
				"blitz":         {},
				"rapid":         {},
				"classical":     {},
				"ultraBullet":   {},
				"chess960":      {},
				"crazyhouse":    {},
				"antichess":     {},
				"atomic":        {},
				"horde":         {},
				"kingOfTheHill": {},
				"racingKings":   {},
				"threeCheck":    {},
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
			lapi.endpoint.topAllPlayers = server.URL

			users, err := lapi.GetAllTop()

			assert.Equal(tt.want, users)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_GetTop(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		Category string
		Number   int
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
			name: "Get top in category",
			args: args{
				Category: "bullet",
				Number:   3,
			},
			params: params{
				query:       "user1,user2",
				requestType: http.MethodGet,
				requestBody: "",
				response: `{
					"users": [
						{
							"id": "bahadirozen",
							"username": "BahadirOzen"
						},
						{
							"id": "penguingim1",
							"username": "penguingim1"
						},
						{
							"id": "night-king96",
							"username": "Night-King96"
						}
				]}`,
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
				{
					ID:       "night-king96",
					Username: "Night-King96",
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
			lapi.endpoint.topPlayers = server.URL + "/%d/%s"

			users, err := lapi.GetTop(tt.args.Category, tt.args.Number)

			assert.Equal(tt.want, users)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_GetUser(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		Username string
	}
	type params struct {
		path        string
		requestType string
		requestBody string
		response    string
	}
	tests := []struct {
		name    string
		args    args
		params  params
		want    User
		wantErr error
	}{
		{
			name: "Get user data",

			params: params{
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
			lapi.endpoint.userProfile = server.URL + "/%s"

			user, err := lapi.GetUser(tt.args.Username)

			assert.Equal(tt.want, *user)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_GetUserRatingHistory(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		Username string
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
		want    map[string][]DailyRating
		wantErr error
	}{
		{
			name: "Get rating history",
			args: args{
				Username: "user1",
			},
			params: params{
				requestType: http.MethodGet,
				requestBody: "",
				response: `[
					{
						"name": "Bullet",
						"points": [
							[
								2011,
								0,
								8,
								1472
							]
						]
					},
					{
						"name": "Blitz",
						"points": [
							[
								2011,
								7,
								29,
								1332
							]
						]
					}]`,
			},
			want: map[string][]DailyRating{
				"Bullet": {
					{2011, 0, 8, 1472},
				},
				"Blitz": {
					{2011, 7, 29, 1332},
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
			lapi.endpoint.userRatingHistory = server.URL + "/%s"

			users, err := lapi.GetUserRatingHistory(tt.args.Username)

			assert.Equal(tt.want, users)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_GetUsesByID(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		IDs []string
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
			name: "Get users",
			args: args{
				IDs: []string{"bahadirozen", "penguingim1"},
			},
			params: params{
				requestType: http.MethodPost,
				requestBody: "bahadirozen,penguingim1",
				response: `[
						{
							"id": "bahadirozen",
							"username": "BahadirOzen"
						},
						{
							"id": "penguingim1",
							"username": "penguingim1"
						}
					]`,
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
			lapi.endpoint.userData = server.URL

			users, err := lapi.GetUsesByID(tt.args.IDs...)

			assert.Equal(tt.want, users)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_GetTeamMembers(t *testing.T) {
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
			name: "Get team members",
			args: args{
				ID: "team",
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
			lapi.endpoint.teamMembers = server.URL + "/%s"

			uchan, _, err := lapi.GetTeamMembers(tt.args.ID)

			var users []User

			for u := range uchan {
				users = append(users, u)
			}

			assert.Equal(tt.want, users)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_GetLiveStreamers(t *testing.T) {
	assert := assert.New(t)

	type params struct {
		requestType string
		requestBody string
		response    string
		query       string
	}
	tests := []struct {
		name    string
		params  params
		want    []User
		wantErr error
	}{
		{
			name: "Get live streamers",
			params: params{
				requestType: http.MethodGet,
				response: `[
					{
						"id": "bahadirozen",
						"username": "BahadirOzen"
					},
					{
						"id": "penguingim1",
						"username": "penguingim1"
					}
				]`,
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
			lapi.endpoint.userLiveStreaming = server.URL

			users, err := lapi.GetLiveStreamers()

			assert.Equal(tt.want, users)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_GetCrosstable(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		userA string
		userB string
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
		want    UserCrosstable
		wantErr error
	}{
		{
			name: "Get crosstable",
			args: args{"userA", "userB"},
			params: params{
				requestType: http.MethodGet,
				response: `{
						"users": 
						{
							"neio": 201.5,
							"thibault": 144.5
						},
						"nbGames": 346
					}`,
			},
			want: UserCrosstable{
				Users: map[string]float32{
					"neio":     201.5,
					"thibault": 144.5,
				},
				NumberGames: 346,
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
			lapi.endpoint.userCrosstable = server.URL + "/%s/%s"

			users, err := lapi.GetCrosstable(tt.args.userA, tt.args.userB)

			assert.Equal(tt.want, *users)
			assert.Equal(tt.wantErr, err)

			users, err = lapi.GetCrosstable(tt.args.userB, tt.args.userA)

			assert.Equal(tt.want, *users)
			assert.Equal(tt.wantErr, err)
		})
	}
}
