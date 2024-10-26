package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nba-sql/nag"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mitchellh/mapstructure"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	playerSearchCmd.Flags().StringP("search", "s", "", "search term")
	playerCmd.Flags().StringVarP(&playerNameFlag, "name", "n", "", "player name")
	playerCmd.AddCommand(playerSearchCmd)
	rootCmd.AddCommand(playerCmd)
}

var playerCmd = &cobra.Command{
	Use:     "players",
	Short:   "Search for players",
	Long:    "",
	Aliases: []string{"p", "player"},
	Run: func(cmd *cobra.Command, args []string) {
		players, err := queryPlayers(DBs["player"])
		if err != nil {
			log.Fatal(err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "Name", "Active"})

		for _, p := range players {
			if fuzzy.MatchFold(playerNameFlag, p.FullName) { // Assuming 'name' is the cobra flag value
				t.AppendRows([]table.Row{
					{p.ID, p.FullName, p.IsActive == 1},
				})
				t.AppendSeparator()
			}
		}
		t.Render()
	},
}

func queryPlayers(db *sql.DB) ([]Player, error) {
	rows, err := db.Query("SELECT * FROM player")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player
		err := rows.Scan(
			&p.ID,
			&p.FullName,
			&p.FirstName,
			&p.LastName,
			&p.IsActive,
		)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}

var playerSearchCmd = &cobra.Command{
	Use:     "search <search terms...>",
	Short:   "Search players for stats",
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		searchTerm := strings.Join(args, " ")
		fmt.Println("Searching for:", searchTerm)

		players := nag.NewCommonAllPlayers()
		err := players.Get()
		if err != nil {
			panic(err)
		}

		if players.Response == nil {
			panic("no response")
		}

		// var result []Player
		var result CommonAllPlayersResponse
		mapstructure.Decode(nag.Map(*players.Response), &result)
		var id int
		// playerMap := map[string]any{}
		for _, data := range result.CommonAllPlayers {
			if fuzzy.MatchFold(searchTerm, data.DISPLAY_FIRST_LAST) {
				fmt.Printf("%#v\n", data)
				id = data.PERSON_ID
				break
			}
			// playerMap[fmt.Sprintf("%d", data.PERSON_ID)] = data.DISPLAY_FIRST_LAST
			// playerMap[data.DISPLAY_FIRST_LAST] = data.PERSON_ID
		}

		playerInfo := nag.NewCommonPlayerInfo(fmt.Sprintf("%d", id))
		err = playerInfo.Get()
		if err != nil {
			panic(err)
		}

		if playerInfo.Response == nil {
			panic("no response")
		}

		p := nag.Map(*playerInfo.Response)

		var playerInfoResult map[string]any
		err = mapstructure.Decode(p, &playerInfoResult)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("%#v\n", p)

		playerGameLog := nag.NewPlayerGameLog(strconv.Itoa(id))
		err = playerGameLog.Get()
		if err != nil {
			log.Fatal(err)
		}

		d := nag.Map(*playerGameLog.Response)

		// fmt.Printf("%#v\n", playerGameLog.Response)

		var playerGameLogResult map[string]any
		err = mapstructure.Decode(d, &playerGameLogResult)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("%#v\n", playerGameLogResult["PlayerGameLog"])

		game_logs := []GameLogEntry{}
		for _, d := range playerGameLogResult["PlayerGameLog"].([]map[string]any) {
			entryJSON, _ := json.Marshal(d)
			var e GameLogEntry
			json.Unmarshal(entryJSON, &e)
			game_logs = append(game_logs, e)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"#",
			"PLAYER_ID",
			"GAME_ID",
			"MATCHUP",
			"GAME_DATE",
			"WL",
			"MIN",
			"PTS",
			"FGM",
			"FGA",
			"FG%",
			"3PM",
			"3PA",
			"3P%",
			"FTM",
			"FTA",
			"FT%",
			"OREB",
			"DREB",
			"REB",
			"AST",
			"STL",
			"BLK",
			"TOV",
			"PF",
			"+/-",
		})
		rows := [][]string{}
		for i, g := range game_logs {
			// fmt.Printf("%#v\n", g)
			row := []string{
				strconv.Itoa(len(game_logs) - i),
				strconv.Itoa(g.Player_ID),
				g.Game_ID,
				g.MATCHUP,
				g.GAME_DATE,
				g.WL,
				strconv.Itoa(g.MIN),
				strconv.Itoa(g.PTS),
				strconv.Itoa(g.FGM),
				strconv.Itoa(g.FGA),
				strconv.FormatFloat(g.FG_PCT, 'f', 3, 64),
				strconv.Itoa(g.FG3M),
				strconv.Itoa(g.FG3A),
				strconv.FormatFloat(g.FG3_PCT, 'f', 3, 64),
				strconv.Itoa(g.FTM),
				strconv.Itoa(g.FTA),
				strconv.FormatFloat(g.FT_PCT, 'f', 3, 64),
				strconv.Itoa(g.OREB),
				strconv.Itoa(g.DREB),
				strconv.Itoa(g.REB),
				strconv.Itoa(g.AST),
				strconv.Itoa(g.STL),
				strconv.Itoa(g.BLK),
				strconv.Itoa(g.TOV),
				strconv.Itoa(g.PF),
				strconv.Itoa(g.PLUS_MINUS),
			}
			rows = append(rows, row)
			colors := make([]tablewriter.Colors, len(row))
			switch g.WL {
			case "W":
				for j := range colors {
					colors[j] = tablewriter.Colors{tablewriter.FgGreenColor}
				}
			case "L":
				for j := range colors {
					colors[j] = tablewriter.Colors{tablewriter.FgRedColor}
				}
			}
			table.Rich(row, colors)
		}
		table.Render()
		// fmt.Printf("%#v\n", playerGameLog.Response.ResultSets.)
		// for _, r := range d.ResultSets {
		// 	// fmt.Printf("%#v\n", r.Headers)
		// 	fmt.Printf("%#v\n", r.RowSet)
		// }

		// repr.Println(playerInfoResult)

		// fmt.Printf("%#v\n", playerInfoResult)

		// if results := fuzzy.RankFind(searchTerm, maps.Keys(playerMap)); len(results) > 0 {
		// 	fmt.Printf("%#v\n", results)
		// }
	},
}
