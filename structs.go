package main

// Structs to keep track of the order of the responses
type PlayersResponse struct {
	Index   int
	Players []Player
}
type PositionsResponse struct {
	Index     int
	RosterMap map[string]string
}

// Struct for how necessary variables are passed to API
type RosterMeta struct {
	LeagueId int    `json:"league_id"`
	EspnS2   string `json:"espn_s2"`
	Swid     string `json:"swid"`
	TeamName string `json:"team_name"`
	Year     int    `json:"year"`
}

// Struct for how to contruct Players using the returned player data
type Player struct {
	Name           string   `json:"name"`
	AvgPoints      float64  `json:"avg_points"`
	Team           string   `json:"team"`
	InjuryStatus   string   `json:"injury_status"`
	ValidPositions []string `json:"valid_positions"`
}

// Struct for JSON schedule file that is used to get days a player is playing
type GameSchedule struct {
	StartDate string                `json:"startDate"`
	EndDate   string                `json:"endDate"`
	GameSpan  int                   `json:"gameSpan"`
	Games     map[string][]int      `json:"games"`
}

// Struct for keeping track of state across recursive function calls to allow for early exit
type FitPlayersContext struct {
	BestLineup map[string]string
	TopScore   int
	MaxScore   int
	EarlyExit  bool
}

// Struct for chromosome for genetic algorithm
type Chromosome struct {
	Genes 	     	 []Gene
	FitnessScore	 int
	TotalAcquisitions int
	CumProbTracker 	 float64
}

// Struct for gene for genetic algorithm
type Gene struct {
	Roster  	 map[string]Player
	NewPlayers 	 map[string]Player
	Day     	 int
	Acquisitions int
}

// Struct to keep track of dropped players during the genetic algorithm
type DroppedPlayer struct {
	Player 	  Player
	Countdown int
}

// Struct to allow for finding the position of the lowest scoring player in a group
type PlayerScore struct {
	Player     Player
	AvgPoints  float64
	Position   string
}