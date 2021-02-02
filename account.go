package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// User struct represents user account in lichess.org
// Some fileds may not be initialized
type User struct {
	ID              string                 `json:"id"`
	Username        string                 `json:"username"`
	Online          bool                   `json:"online"`
	Performances    map[string]Performance `json:"perfs"`
	CreatedAt       int64                  `json:"createdAt"`
	Disabled        bool                   `json:"disabled"`
	TOSViolation    bool                   `json:"tosViolation"`
	Booster         bool                   `json:"booster"`
	Profile         Profile                `json:"profile"`
	SeenAt          int64                  `json:"seenAt"`
	Patron          bool                   `json:"patron"`
	PlayTime        PlayTime               `json:"playTime"`
	Language        string                 `json:"language"`
	Title           string                 `json:"title"`
	URL             string                 `json:"url"`
	Playing         string                 `json:"playing"`
	NumberFollowing int                    `json:"nbFollowing"`
	NumberFollowers int                    `json:"nbFollowers"`
	CompletionRate  int                    `json:"completionRate"`
	Count           StatsCount             `json:"count"`
	Streaming       bool                   `json:"streaming"`
	Followable      bool                   `json:"followable"`
	Following       bool                   `json:"following"`
	Blocking        bool                   `json:"blocking"`
	FollowsYou      bool                   `json:"followsYou"`
}

// Performance struct stores user performance in one category
type Performance struct {
	Games           int  `json:"games"`
	Rating          int  `json:"rating"`
	RatingDeviation int  `json:"rd"`
	Prog            int  `json:"prog"` // IDK what this field means
	Prov            bool `json:"prov"` // IDK what this field means
}

// Profile struct adds additional information about user
type Profile struct {
	Country    string `json:"country"`
	Location   string `json:"location"`
	Bio        string `json:"bio"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	FideRating int    `json:"fideRating"`
	UscfRating int    `json:"uscfRating"`
	EcfRating  int    `json:"ecfRating"`
	Links      string `json:"links"`
}

// PlayTime stores information about time played in minites
type PlayTime struct {
	Total int `json:"total"`
	TV    int `json:"tv"`
}

// StatsCount stores game results
type StatsCount struct {
	All      int `json:"all"`
	Rated    int `json:"rated"`
	AI       int `json:"ai"`
	Draw     int `json:"draw"`
	DrawH    int `json:"drawH"` // IDK what this field means
	Loss     int `json:"loss"`
	LossH    int `json:"lossH"` // IDK what this field means
	Win      int `json:"win"`
	WinH     int `json:"winH"` // IDK what this field means
	Bookmark int `json:"bookmark"`
	Playing  int `json:"playing"`
	Import   int `json:"import"`
	Me       int `json:"me"`
}

// Preferences stores user's UI preferences
type Preferences struct {
	Dark          bool   `json:"dark"`
	Transparent   bool   `json:"transp"`
	BgImg         string `json:"bgImg"`
	Is3d          bool   `json:"is3d"`
	Theme         string `json:"theme"`
	PieceSet      string `json:"pieceSet"`
	Theme3d       string `json:"theme3d"`
	PieceSet3d    string `json:"pieceSet3d"`
	SoundSet      string `json:"soundSet"`
	BlindFold     int    `json:"blindfold"`
	AutoQueen     int    `json:"autoQueen"`
	AutoThreeFold int    `json:"autoThreefold"`
	Takeback      int    `json:"takeback"`
	Moretime      int    `json:"moretime"`
	ClockTenths   int    `json:"clockTenths"`
	ClockBar      bool   `json:"clockBar"`
	ClockSound    bool   `json:"clockSound"`
	Premove       bool   `json:"premove"`
	Animation     int    `json:"animation"`
	Captured      bool   `json:"captured"`
	Follow        bool   `json:"follow"`
	Highlight     bool   `json:"highlight"`
	Destination   bool   `json:"destination"`
	Coords        int    `json:"coords"`
	Replay        int    `json:"replay"`
	Challenge     int    `json:"challenge"`
	Message       int    `json:"message"`
	CoordColor    int    `json:"coordColor"`
	SubmitMove    int    `json:"submitMove"`
	ConfirmResign int    `json:"confirmResign"`
	InsightShare  int    `json:"insightShare"`
	KeyboardMove  int    `json:"keyboardMove"`
	Zen           int    `json:"zen"`
	MoveEvent     int    `json:"moveEvent"`
	RookCastle    int    `json:"rookCastle"`
}

// GetMyProfile returns information about logged user
func (l *LichessAPI) GetMyProfile() (*User, error) {
	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    l.endpoint.accountProfile,
	}

	resp, err := l.request(params)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetMyEmail returns user's email
func (l *LichessAPI) GetMyEmail() (string, error) {
	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    l.endpoint.accountEmail,
	}

	resp, err := l.request(params)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	type email struct {
		Email string `json:"email"`
	}

	var e email

	err = json.NewDecoder(resp.Body).Decode(&e)
	if err != nil {
		return "", err
	}

	return e.Email, nil
}

// GetMyPreferences returns user's preferences
func (l *LichessAPI) GetMyPreferences() (*Preferences, error) {
	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    l.endpoint.accountPreferences,
	}

	resp, err := l.request(params)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	type prefsResp struct {
		Prefs Preferences `json:"prefs"`
	}

	var prefs prefsResp

	err = json.NewDecoder(resp.Body).Decode(&prefs)
	if err != nil {
		return nil, err
	}

	return &prefs.Prefs, nil
}

// GetMyKidModeStatus returns user's kid mode status
func (l *LichessAPI) GetMyKidModeStatus() (bool, error) {
	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    l.endpoint.accountKidModeStatus,
	}

	resp, err := l.request(params)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	type kidStatus struct {
		Kid bool `json:"kid"`
	}

	var kid kidStatus

	err = json.NewDecoder(resp.Body).Decode(&kid)
	if err != nil {
		return false, err
	}

	return kid.Kid, nil
}

// SetMyKidModeStatus sets user's kid mode status.
// Returns true on success
func (l *LichessAPI) SetMyKidModeStatus(newStatus bool) (bool, error) {
	params := &reqParams{
		requestType: http.MethodPost,
		endpoint:    l.endpoint.accountKidModeStatus,
		data:        []byte(fmt.Sprintf("v=%v", newStatus)),
	}

	resp, err := l.request(params)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	type result struct {
		Ok bool `json:"ok"`
	}

	var res result
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return false, err
	}

	return res.Ok, nil
}
