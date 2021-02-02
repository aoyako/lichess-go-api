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
}

func newServiceEndpoint() *serviceEndpoint {
	return &serviceEndpoint{
		accountProfile:       "https://lichess.org/api/account",
		accountEmail:         "https://lichess.org/api/account/email",
		accountPreferences:   "https://lichess.org/api/account/preferences",
		accountKidModeStatus: "https://lichess.org/api/account/kid",

		userStatus:    "https://lichess.org/api/users/status",
		topAllPlayers: "https://lichess.org/player",
		topPlayers:    "https://lichess.org/player/top/%d/%s",
		userProfile:   "https://lichess.org/api/user/%s",
	}
}
