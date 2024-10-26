package cmd

import (
	"fmt"
	"strings"

	"github.com/k0kubun/pp/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(PlayByPlayCmd)
}

var PlayByPlayCmd = &cobra.Command{
	Use:     "playbyplay",
	Aliases: []string{"pbp", "play"},
	Short:   "",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		playbyplay()
	},
}

func playbyplay() {
	db := DBs["play-by-play"]
	rows, _ := db.Query("SELECT * FROM play_by_play LIMIT 10")
	defer rows.Close()

	rowCount := 0

	// Iterate through the rows
	for rows.Next() {
		var pbp PlayByPlay

		// Scan the row directly into our struct
		if err := rows.Scan(
			&pbp.GameID,
			&pbp.EventNum,
			&pbp.EventMsgType,
			&pbp.EventMsgActionType,
			&pbp.Period,
			&pbp.WCTimeString,
			&pbp.PCTimeString,
			&pbp.HomeDescription,
			&pbp.NeutralDescription,
			&pbp.VisitorDescription,
			&pbp.Score,
			&pbp.ScoreMargin,
			&pbp.Person1Type,
			&pbp.Player1ID,
			&pbp.Player1Name,
			&pbp.Player1TeamID,
			&pbp.Player1TeamCity,
			&pbp.Player1TeamNickname,
			&pbp.Player1TeamAbbrev,
			&pbp.Person2Type,
			&pbp.Player2ID,
			&pbp.Player2Name,
			&pbp.Player2TeamID,
			&pbp.Player2TeamCity,
			&pbp.Player2TeamNickname,
			&pbp.Player2TeamAbbrev,
			&pbp.Person3Type,
			&pbp.Player3ID,
			&pbp.Player3Name,
			&pbp.Player3TeamID,
			&pbp.Player3TeamCity,
			&pbp.Player3TeamNickname,
			&pbp.Player3TeamAbbrev,
			&pbp.VideoAvailableFlag,
		); err != nil {
			log.Fatal(err)
		}

		rowCount++
		pp.Print(pbp)
		println()
	}

	// Print footer
	fmt.Printf("\nTotal Plays: %d\n", rowCount)
	fmt.Println(strings.Repeat("=", 40))

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
