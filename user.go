package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// DailyRating stores rating during ont day
type DailyRating struct {
	Year   int
	Month  int
	Day    int
	Rating int
}

// Activity stores user activity on lichess.org
type Activity struct {
	Interval            Interval                `json:"interval"`
	Games               map[string]GameActivity `json:"games"`
	Puzzles             GameActivity            `json:"puzzles"`
	Tournaments         TournamentActivity      `json:"tournaments"`
	Practices           []Practice              `json:"practice"`
	CorrespondenceMoves CorrespondenceMoves     `json:"correspondenceMoves"`
	CorrespondenceEnds  CorrespondenceEnds      `json:"correspondenceEnds"`
	Teams               []TeamActivity          `json:"teams"`
	Posts               []Topic                 `json:"posts"`
}

// Topic represents forum topic
type Topic struct {
	TopicURL  string `json:"topicUrl"`
	TopicName string `json:"topicName"`
	Posts     []Post `json:"posts"`
}

// Post represents forum post
type Post struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}

// TeamActivity stores user activity in team
type TeamActivity struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

// CorrespondenceMoves stores corresponcence active games
type CorrespondenceMoves struct {
	Nb    int            `json:"nb"`
	Games []GameByPlayer `json:"games"`
}

// CorrespondenceEnds stores corresponcence ended games
type CorrespondenceEnds struct {
	Score GameActivity   `json:"score"`
	Games []GameByPlayer `json:"games"`
}

// Follows stores user ids that have followed or unfollowed user
type Follows struct {
	In  []string
	Out []string
}

// UnmarshalJSON for Follows struct
func (f *Follows) UnmarshalJSON(data []byte) error {
	type follow struct {
		IDs []string `json:"ids"`
	}
	type IOFollow struct {
		In  follow `json:"in"`
		Out follow `json:"out"`
	}

	var v IOFollow
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	f.In = v.In.IDs
	f.Out = v.Out.IDs

	return nil
}

// GameActivity stores results of finished games
type GameActivity struct {
	Win  int          `json:"win"`
	Loss int          `json:"loss"`
	Draw int          `json:"draw"`
	Rp   RatingChange `json:"rp"`
}

// RatingChange stores rating difference
type RatingChange struct {
	Before int `json:"before"`
	After  int `json:"after"`
}

// TournamentActivity stores user's tournament activity
type TournamentActivity struct {
	Nb   int               `json:"nb"`
	Best []TournamentScore `json:"best"`
}

// TournamentScore stores user's tournament score
type TournamentScore struct {
	Tournament  ArenaTournament `json:"tournament"`
	NbGames     int             `json:"nbGames"`
	Score       int             `json:"score"`
	Rank        int             `json:"rank"`
	RankPercent int             `json:"rankPercent"`
}

// Practice represents lesson in lichess
type Practice struct {
	URL         string `json:"url"`
	Name        string `json:"name"`
	NbPositions int    `json:"nbPositions"`
}

// ArenaTournament stores tournament information
type ArenaTournament struct {
	ID             string         `json:"id"`
	CreatedBy      string         `json:"createdBy"`
	System         string         `json:"system"`
	Minutes        int            `json:"minutes"`
	Clock          Clock          `json:"clock"`
	Rated          bool           `json:"rated"`
	FullName       string         `json:"fullName"`
	Name           string         `json:"name"`
	NbPlayers      int            `json:"nbPlayers"`
	Variant        Variant        `json:"variant"`
	StartsAt       int64          `json:"startsAt"`
	FinishesAt     int64          `json:"finishesAt"`
	Status         int            `json:"status"`
	Perf           TournamentPerf `json:"perf"`
	SecondsToStart int            `json:"secondsToStart"`
	HasMaxRating   bool           `json:"hasMaxRating"`
	Private        bool           `json:"private"`
	Position       Opening        `json:"position"`
	Schedule       Schedule       `json:"schedule"`
	Winner         User           `json:"winner"`
}

// Variant represents game type
type Variant struct {
	Key   string `json:"key"`
	Short string `json:"short"`
	Name  string `json:"name"`
}

// TournamentPerf stores tournament parameters
type TournamentPerf struct {
	Icon     string `json:"icon"`
	Key      string `json:"key"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

// Schedule stores tournament schedule
type Schedule struct {
	Freq  string `json:"freq"`
	Speed string `json:"speed"`
}

// Opening represents opening position
type Opening struct {
	Eco      string `json:"eco"`
	Fen      string `json:"fen"`
	Name     string `json:"name"`
	Ply      int    `json:"ply"`
	WikiPath string `json:"wikiPath"`
}

// GameByPlayer represents game played by user
type GameByPlayer struct {
	ID       string       `json:"id"`
	FullID   string       `json:"fullId"`
	GameID   string       `json:"gameID"`
	Color    string       `json:"color"`
	URL      string       `json:"url"`
	Variant  string       `json:"variant"`
	Speed    string       `json:"speed"`
	Perf     string       `json:"perf"`
	Rated    bool         `json:"rated"`
	Opponent GameOpponent `json:"opponent"`
	Fen      string       `json:"fen"`
	IsMyTurn bool         `json:"isMyTurn"`
}

// GameOpponent stores basic info about user's opponent
type GameOpponent struct {
	ID     string `json:"id"`
	User   string `json:"user"`
	Rating int    `json:"rating"`
}

// Interval represents time interval between activities
type Interval struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// Clock represents chess clock
type Clock struct {
	Initial   int `json:"initial"`
	Increment int `json:"increment"`
	TotalTime int `json:"totalTime"`
}

// UserCrosstable stores score of users matching against each other
type UserCrosstable struct {
	Users       map[string]float32 `json:"users"`
	NumberGames int                `json:"nbGames"`
	Matchup     *UserCrosstable    `json:"matchup"`
}

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
		sb.WriteString(",")
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

// GetUserRatingHistory returns rating history of given user.
// Format: map "category" -> array of ratings
func (l *LichessAPI) GetUserRatingHistory(username string) (map[string][]DailyRating, error) {
	type ratingHistory struct {
		Name   string  `json:"name"`
		Points [][]int `json:"points"`
	}

	var history []ratingHistory

	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    fmt.Sprintf(l.endpoint.userRatingHistory, username),
	}

	resp, err := l.request(params)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&history)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]DailyRating)
	for _, helem := range history {
		ratings := make([]DailyRating, len(helem.Points))

		for _, point := range helem.Points {
			ratings = append(ratings, DailyRating{point[0], point[1], point[2], point[3]})
		}

		result[helem.Name] = ratings
	}

	return result, nil
}

// GetUserActivity returns user activity
func (l *LichessAPI) GetUserActivity(username string) ([]Activity, error) {
	var activity []Activity

	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    fmt.Sprintf(l.endpoint.userActivity, username),
	}

	resp, err := l.request(params)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&activity)
	if err != nil {
		return nil, err
	}

	return activity, nil
}

// GetUsesByID returns users by their ids
func (l *LichessAPI) GetUsesByID(ids ...string) ([]User, error) {
	var users []User

	if len(ids) == 0 {
		return users, nil
	}
	if len(ids) > 300 {
		return users, errors.New("Too many requested profiles, max is 300")
	}

	var sb strings.Builder

	for pos := range ids[:len(ids)-1] {
		sb.WriteString(ids[pos])
		sb.WriteString(",")
	}
	sb.WriteString(ids[len(ids)-1])

	params := &reqParams{
		requestType: http.MethodPost,
		endpoint:    l.endpoint.userData,
		data:        []byte(sb.String()),
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

// GetTeamMembers returns team members.
// Use channel to get streamed values.
// Call returned function to stop receiving
func (l *LichessAPI) GetTeamMembers(id string) (chan User, func(), error) {
	users := make(chan User, 10)
	finishContext, cancel := context.WithCancel(context.Background())

	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    fmt.Sprintf(l.endpoint.teamMembers, id),
	}

	resp, err := l.request(params)
	if err != nil {
		return nil, cancel, err
	}

	reader := bufio.NewReader(resp.Body)

	go func() {
		defer resp.Body.Close()

		for {
			select {
			case <-finishContext.Done():
				return
			default:
				message, err := reader.ReadBytes('\n')
				if err != nil {
					close(users)
					return
				}
				if len(message) != 0 {
					var user User
					err = json.NewDecoder(bytes.NewReader(message)).Decode(&user)
					if err != nil {
						close(users)
						return
					}

					users <- user
				}
			}
		}
	}()

	return users, cancel, nil
}

// GetLiveStreamers returs current live streaming users
func (l *LichessAPI) GetLiveStreamers() ([]User, error) {
	var users []User

	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    l.endpoint.userLiveStreaming,
	}

	resp, err := l.request(params)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetCrosstable returs crosstable of given users
func (l *LichessAPI) GetCrosstable(usernameA, usernameB string) (*UserCrosstable, error) {
	var result UserCrosstable

	params := &reqParams{
		requestType: http.MethodGet,
		endpoint:    fmt.Sprintf(l.endpoint.userCrosstable, usernameA, usernameB),
	}

	resp, err := l.request(params)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
