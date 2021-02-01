package main

type serviceEndpoint struct {
	accountProfile       string
	accountEmail         string
	accountPreferences   string
	accountKidModeStatus string
}

func newServiceEndpoint() *serviceEndpoint {
	return &serviceEndpoint{
		accountProfile:       "https://lichess.org/api/account",
		accountEmail:         "https://lichess.org/api/account/email",
		accountPreferences:   "https://lichess.org/api/account/preferences",
		accountKidModeStatus: "https://lichess.org/api/account/kid",
	}
}
