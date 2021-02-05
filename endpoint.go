package main

type serviceEndpoint struct {
	accountProfile       string
	accountEmail         string
	accountPreferences   string
	accountKidModeStatus string
	userStatus           string
	topAllPlayers        string
	topPlayers           string
	userProfile          string
	userRatingHistory    string
	userActivity         string
	userData             string
	teamMembers          string
	userLiveStreaming    string
	userCrosstable       string
	relationFollowing    string
	relationFollowers    string
}

func newServiceEndpoint() *serviceEndpoint {
	return &serviceEndpoint{
		accountProfile:       "https://lichess.org/api/account",
		accountEmail:         "https://lichess.org/api/account/email",
		accountPreferences:   "https://lichess.org/api/account/preferences",
		accountKidModeStatus: "https://lichess.org/api/account/kid",

		userStatus:        "https://lichess.org/api/users/status",
		topAllPlayers:     "https://lichess.org/player",
		topPlayers:        "https://lichess.org/player/top/%d/%s",
		userProfile:       "https://lichess.org/api/user/%s",
		userRatingHistory: "https://lichess.org/api/user/%s/rating-history",
		userActivity:      "https://lichess.org/api/user/%s/activity",
		userData:          "https://lichess.org/api/users",
		teamMembers:       "https://lichess.org/api/team/%s/users",
		userLiveStreaming: "https://lichess.org/streamer/live",
		userCrosstable:    "https://lichess.org/api/crosstable/%s/%s",

		relationFollowing: "https://lichess.org/api/user/%s/following",
		relationFollowers: "https://lichess.org/api/user/%s/followers",
	}
}
