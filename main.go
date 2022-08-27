package main

import (
	"log"
	"os"
	"time"

	"goapiproject/handlers"
	"goapiproject/helpers"

	"github.com/asdine/storm/v3"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	bolt "go.etcd.io/bbolt"
)

var apiPort string

type Achievement struct {
	ID          int `storm:"id,increment"`
	Name        string
	Description string
}

func main() {
	app := &cli.App{
		Name:   "apiserver",
		Usage:  "A very basic REST API with basic database backing",
		Action: apiServer,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"api"},
				Usage:       "assign `PORT` to send API calls to",
				Destination: &apiPort,
				Value:       "8080",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// apiServer runs the API server using gin, bbolt, and storm
func apiServer(c *cli.Context) error {
	serverIP, _ := helpers.GetServersOwnIP() // retrieve own IP for logging purposes
	log.Printf("Server started\n\tAPI accessible on %v:%s\n", serverIP, apiPort)

	// Open the api_data.db data file in the current directory.
	// DB file will be created if it doesn't exist. 1s timeout prevents indefinite wait
	db, err := storm.Open("api_data.db", storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // close database connection to release file lock

	err = helpers.ResetAndFillDB(db) // reset the DB and put some initial data in it
	if err != nil {
		return err
	}

	gin.SetMode(gin.ReleaseMode) // set automatic logging level
	router := gin.Default()      // use default gin instance with logger and recovery middleware preattached

	// set up routers and their respective handler functions
	router.GET("/online", handlers.Online())
	router.GET("/achievement/all", handlers.GetAllAchievements(db))
	router.GET("/achievement/one", handlers.GetOneAchievement(db))
	router.POST("/achievement/create", handlers.CreateANewAchievement(db))
	router.PUT("/achievement/update", handlers.UpdateExistingAchievement(db))
	router.DELETE("/achievement/delete", handlers.DeleteExistingAchievement(db))
	router.PUT("/achievement/reset/really", handlers.ResetDatabase(db))

	// listen and serve on <instance localhost/IP>:<port>
	// accessible through the wm's eth0 IP in WSL
	err = router.Run(":" + apiPort)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
