package cmd

import (
	"database/sql"

	"github.com/tursodatabase/go-libsql"
)

type GameLogEntry struct {
	AST             int     `json:"AST"`
	BLK             int     `json:"BLK"`
	DREB            int     `json:"DREB"`
	FG3A            int     `json:"FG3A"`
	FG3M            int     `json:"FG3M"`
	FG3_PCT         float64 `json:"FG3_PCT"`
	FGA             int     `json:"FGA"`
	FGM             int     `json:"FGM"`
	FG_PCT          float64 `json:"FG_PCT"`
	FTA             int     `json:"FTA"`
	FTM             int     `json:"FTM"`
	FT_PCT          float64 `json:"FT_PCT"`
	GAME_DATE       string  `json:"GAME_DATE"`
	Game_ID         string  `json:"Game_ID"`
	MATCHUP         string  `json:"MATCHUP"`
	MIN             int     `json:"MIN"`
	OREB            int     `json:"OREB"`
	PF              int     `json:"PF"`
	PLUS_MINUS      int     `json:"PLUS_MINUS"`
	PTS             int     `json:"PTS"`
	Player_ID       int     `json:"Player_ID"`
	REB             int     `json:"REB"`
	SEASON_ID       string  `json:"SEASON_ID"`
	STL             int     `json:"STL"`
	TOV             int     `json:"TOV"`
	VIDEO_AVAILABLE int     `json:"VIDEO_AVAILABLE"`
	WL              string  `json:"WL"`
}

type Player2 struct {
	DISPLAY_FIRST_LAST        string `json:"DISPLAY_FIRST_LAST"`
	DISPLAY_LAST_COMMA_FIRST  string `json:"DISPLAY_LAST_COMMA_FIRST"`
	FROM_YEAR                 int    `json:"FROM_YEAR"`
	GAMES_PLAYED_FLAG         string `json:"GAMES_PLAYED_FLAG"`
	OTHERLEAGUE_EXPERIENCE_CH string `json:"OTHERLEAGUE_EXPERIENCE_CH"`
	PERSON_ID                 int    `json:"PERSON_ID"`
	PLAYERCODE                string `json:"PLAYERCODE"`
	PLAYER_SLUG               string `json:"PLAYER_SLUG"`
	ROSTERSTATUS              int    `json:"ROSTERSTATUS"`
	TEAM_ABBREVIATION         string `json:"TEAM_ABBREVIATION"`
	TEAM_CITY                 string `json:"TEAM_CITY"`
	TEAM_CODE                 string `json:"TEAM_CODE"`
	TEAM_ID                   int64  `json:"TEAM_ID"`
	TEAM_NAME                 string `json:"TEAM_NAME"`
	TEAM_SLUG                 string `json:"TEAM_SLUG"`
	TO_YEAR                   int    `json:"TO_YEAR"`
}

type CommonAllPlayersResponse struct {
	CommonAllPlayers []Player2 `json:"CommonAllPlayers"`
}

type PlayByPlay struct {
	GameID              string `db:"game_id"`
	EventNum            string `db:"eventnum"`
	EventMsgType        string `db:"eventmsgtype"`
	EventMsgActionType  string `db:"eventmsgactiontype"`
	Period              string `db:"period"`
	WCTimeString        string `db:"wctimestring"`
	PCTimeString        string `db:"pctimestring"`
	HomeDescription     string `db:"homedescription"`
	NeutralDescription  string `db:"neutraldescription"`
	VisitorDescription  string `db:"visitordescription"`
	Score               string `db:"score"`
	ScoreMargin         string `db:"scoremargin"`
	Person1Type         string `db:"person1type"`
	Player1ID           string `db:"player1_id"`
	Player1Name         string `db:"player1_name"`
	Player1TeamID       string `db:"player1_team_id"`
	Player1TeamCity     string `db:"player1_team_city"`
	Player1TeamNickname string `db:"player1_team_nickname"`
	Player1TeamAbbrev   string `db:"player1_team_abbreviation"`
	Person2Type         string `db:"person2type"`
	Player2ID           string `db:"player2_id"`
	Player2Name         string `db:"player2_name"`
	Player2TeamID       string `db:"player2_team_id"`
	Player2TeamCity     string `db:"player2_team_city"`
	Player2TeamNickname string `db:"player2_team_nickname"`
	Player2TeamAbbrev   string `db:"player2_team_abbreviation"`
	Person3Type         string `db:"person3type"`
	Player3ID           string `db:"player3_id"`
	Player3Name         string `db:"player3_name"`
	Player3TeamID       string `db:"player3_team_id"`
	Player3TeamCity     string `db:"player3_team_city"`
	Player3TeamNickname string `db:"player3_team_nickname"`
	Player3TeamAbbrev   string `db:"player3_team_abbreviation"`
	VideoAvailableFlag  string `db:"video_available_flag"`
}

type DBConnection struct {
	name      string
	db        *sql.DB
	connector *libsql.Connector
}

type DBResult struct {
	name string
	db   *sql.DB
	err  error
}

type Player struct {
	ID        string
	FullName  string
	FirstName string
	LastName  string
	IsActive  int
}
