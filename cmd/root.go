package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/chelnak/ysmrr"
	"github.com/chelnak/ysmrr/pkg/animations"
	"github.com/chelnak/ysmrr/pkg/colors"
	"github.com/joho/godotenv"
	"github.com/shogo82148/go-shuffle"
	log "github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/cobra"
	"github.com/tursodatabase/go-libsql"
)

const (
	suffix = "ksc98.turso.io"
)

var (
	syncFlag        bool
	parallelismFlag int
	playerNameFlag  string
	DBs             = make(map[string]*sql.DB)
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&syncFlag, "sync", "s", false, "sync databases to remote")
	rootCmd.PersistentFlags().IntVarP(&parallelismFlag, "parallel", "p", 20, "number of concurrent replica syncs")
}

var databaseNames = []string{
	"team-info-common",    // 245B
	"team",                // 2.0KB
	"team-history",        // 2.0KB
	"team-details",        // 5.4KB
	"player",              // 166KB
	"draft-combine-stats", // 199KB
	"draft-history",       // 810KB
	"common-player-info",  // 1,006KB
	"officials",           // 2.2MB
	"game-info",           // 2.3MB
	"other-stats",         // 3.3MB
	"game-summary",        // 5.4MB
	"inactive-players",    // 7.0MB
	"line-score",          // 10MB
	"game",                // 19MB
	"play-by-play",        // 2.1GB, remote
}

func load_dotenv() {
	err := godotenv.Load() // load .env
	if err != nil {
		log.Fatal("error loading .env file: ", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "nbacli",
	Short: "NBA Database Engine & CLI",
	Long:  `Query and synchronize NBA data with remote Turso databases`,
	Run:   func(cmd *cobra.Command, args []string) {},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		load_dotenv()
		var (
			authToken = os.Getenv("TURSO_AUTH_TOKEN")
			results   = make(chan DBResult, len(databaseNames))
			sm        ysmrr.SpinnerManager
			p         = pool.New().WithMaxGoroutines(parallelismFlag)
		)

		if syncFlag {
			sm = ysmrr.NewSpinnerManager(
				ysmrr.WithAnimation(animations.Dots),
				ysmrr.WithSpinnerColor(colors.FgHiYellow),
			)
			sm.Start()
		}

		shuffle.Strings(databaseNames)
		for _, dbName := range databaseNames {
			var (
				spinner        *ysmrr.Spinner
				spinnerErrored bool
			)

			if syncFlag {
				spinner = sm.AddSpinner(fmt.Sprintf("syncing db [%s]...", dbName))
			}

			p.Go(func() {
				defer func() {
					if syncFlag && !spinnerErrored {
						spinner.Complete()
					}
				}()

				primaryUrl := fmt.Sprintf("libsql://%s-%s", dbName, suffix)
				dbFilePath := fmt.Sprintf("./db/%s.sql", dbName)

				dbResult := DBResult{
					name: dbName,
				}

				handleErr := func(err error) {
					dbResult.err = err
					results <- dbResult
					if syncFlag {
						spinner.Error()
						spinnerErrored = true
					}
				}

				if dbName == "play-by-play" {
					primaryUrl += fmt.Sprintf("?authToken=%s", authToken)
					db, err := sql.Open("libsql", primaryUrl)
					if err != nil {
						log.Errorf("failed to open db %s: %s\n", primaryUrl, err)
					}
					if err = db.Ping(); err != nil {
						log.Fatalf("error pinging db [%s]: %v", dbName, err)
					}
					dbResult.db = db
					results <- dbResult
					return
				}

				connector, err := libsql.NewEmbeddedReplicaConnector(
					dbFilePath,
					primaryUrl,
					libsql.WithAuthToken(authToken),
					libsql.WithSyncInterval(time.Minute),
				)
				if err != nil {
					handleErr(fmt.Errorf("error creating connector for [%s]: %v", dbName, err))
					return
				}

				// sync db
				if syncFlag {
					if _, err := connector.Sync(); err != nil {
						connector.Close()
						handleErr(fmt.Errorf("error syncing db [%s]: %v", dbName, err))
						return
					}
				}

				dbResult.db = sql.OpenDB(connector)
				results <- dbResult
			})
		}
		p.Wait()
		close(results)
		if syncFlag {
			sm.Stop()
		}

		for dbResult := range results {
			if dbResult.err == nil {
				DBs[dbResult.name] = dbResult.db
			} else {
				log.Errorf("failed connecting to db [%s]: %v", dbResult.name, dbResult.err)
				DBs[dbResult.name] = nil
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
