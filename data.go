package lchess

type Authorization struct {
	Token string
}

type Account struct {
	ID             string          `json:"id"`
	Username       string          `json:"username"`
	Online         bool            `json:"online"`
	Perfs          map[string]Perf `json:"perfs"`
	CreatedAt      int             `json:"createdAt"`
	Disabled       bool            `json:"disabled"`
	TosViolation   bool            `json:"tosViolation"`
	Booster        bool            `json:"booster"`
	Profile        Profile         `json:"profile"`
	SeenAt         int             `json:"seenAt"`
	Patron         bool            `json:"patron"`
	PlayTime       PlayTime        `json:"playTime"`
	Language       string          `json:"language"`
	Title          string          `json:"title"`
	URL            string          `json:"url"`
	Playing        string          `json:"playing"`
	NbFollowing    int             `json:"nbFollowing"`
	NbFollowers    int             `json:"nbFollowers"`
	CompletionRate int             `json:"completionRate"`
	Count          Count           `json:"count"`
	Streaming      bool            `json:"streaming"`
	Followable     bool            `json:"followable"`
	Following      bool            `json:"following"`
	Blocking       bool            `json:"blocking"`
	FollowsYou     bool            `json:"followsYou"`
}

type Perf struct {
	Games  int  `json:"games"`
	Rating int  `json:"rating"`
	Rd     int  `json:"rd"`
	Prog   int  `json:"prog"`
	Prov   bool `json:"prov"`
}

type Profile struct {
	Country    string `json:"counry"`
	Location   string `json:"location"`
	Bio        string `json:"bio"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	FideRating int    `json:"fideRating"`
	UscfRating int    `json:"uscfRating"`
	EcfRating  int    `json:"ecfRating"`
	Links      string `json:"links"`
}

type PlayTime struct {
	Total int `json:"total"`
	Tv    int `json:"tv"`
}

type Count struct {
	All      int `json:"all"`
	Rated    int `json:"rated"`
	Ai       int `json:"ai"`
	Draw     int `json:"draw"`
	DrawH    int `json:"drawH"`
	Loss     int `json:"loss"`
	LossH    int `json:"lossH"`
	Win      int `json:"win"`
	WinH     int `json:"winH"`
	Bookmark int `json:"bookmark"`
	Playing  int `json:"playing"`
	Import   int `json:"import"`
	Me       int `json:"me"`
}

type Prefs struct {
	Dark          bool   `json:"dark"`
	Transp        bool   `json:"transp"`
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

type TotalRatingHistory struct {
	Variants []VariantRatingHistory
}

type VariantRatingHistory struct {
	Name    string          `json:"name"`
	History []RatingHistory `json:"points"`
}

type RatingHistory struct {
	Year   int
	Month  int
	Day    int
	Rating int
}

type Activity struct {
	Interval            Interval                `json:"interval"`
	Games               map[string]GameActivity `json:"games"`
	Puzzles             GameActivity            `json:"puzzles"`
	Tournaments         TournamentActivity      `json:"tournaments"`
	Practices           []Practice              `json:"practice"`
	CorrespondenceMoves CorrespondenceMoves     `json:"correspondenceMoves"`
	CorrespondenceEnds  CorrespondenceEnds      `json:"correspondenceEnds"`
	Teams               []Team                  `json:"teams"`
	Posts               []Topic                 `json:"posts"`
}

type Topic struct {
	TopicURL  string `json:"topicUrl"`
	TopicName string `json:"topicName"`
	Posts     []Post `json:"posts"`
}

type Post struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}

type Team struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type CorrespondenceMoves struct {
	Nb    int    `json:"nb"`
	Games []Game `json:"games"`
}

type CorrespondenceEnds struct {
	Score GameActivity `json:"score"`
	Games []Game       `json:"games"`
}

type Follows struct {
	In  Follow `json:"in"`
	Out Follow `json:"out"`
}

type Follow struct {
	IDs []string `json:"ids"`
}

type Game struct {
	ID       string       `json:"id"`
	Color    string       `json:"color"`
	URL      string       `json:"url"`
	Variant  string       `json:"variant"`
	Speed    string       `json:"speed"`
	Perf     string       `json:"perf"`
	Rated    bool         `json:"rated"`
	Opponent GameOpponent `json:"opponent"`
}

type GameOpponent struct {
	User   string `json:"user"`
	Rating int    `json:"rating"`
}

type Interval struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type GameActivity struct {
	Win  int          `json:"win"`
	Loss int          `json:"loss"`
	Draw int          `json:"draw"`
	Rp   RatingChange `json:"rp"`
}

type RatingChange struct {
	Before int `json:"before"`
	After  int `json:"after"`
}

type TournamentActivity struct {
	Nb   int               `json:"nb"`
	Best []TournamentScore `json:"best"`
}

type TournamentScore struct {
	Tournament  ArenaTournament `json:"tournament"`
	NbGames     int             `json:"nbGames"`
	Score       int             `json:"score"`
	Rank        int             `json:"rank"`
	RankPercent int             `json:"rankPercent"`
}

type ArenaTournament struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Practice struct {
	URL         string `json:"url"`
	Name        string `json:"name"`
	NbPositions int    `json:"nbPositions"`
}
