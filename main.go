package main

import (
	"fmt"
	"log"
	"os"
	"time"

	dbactions "goapiproject/dbactions"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	bolt "go.etcd.io/bbolt"
)

var apiPort string

func main() {

	app := &cli.App{
		Name:   "apiserver",
		Usage:  "API endpoints are accessible on <localhost> or <IP of the machine where it runs, port number is specified as an argument (default: 8080)", // TODO
		Action: apiserv,
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

// TODO do we want to use this somehow, or at all?
type urlDB struct {
	dataBase *bolt.DB
}

// apiserv runs the API server using gin
func apiserv(c *cli.Context) error {
	log.Printf("Server started")
	log.Printf("\tAPI accessible on port: %s\n", apiPort)
	fmt.Println()

	//TODO how to pass into handlerfuncs in reasonable way
	// Open the api_data.db data file in the current directory.
	// DB file will be created if it doesn't exist. 1s timeout prevents indefinite wait
	db, err := bolt.Open("api_data.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // close to release file lock

	err = dbactions.CreateAndFillBuckets(db)
	if err != nil {
		return err
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET( // GET request on /ping endpoint
		"/ping", // (relative) PATH
		// TODO return random message from the database, consider "prefix scans"
		func(c *gin.Context) { // handlerfunc (what to do when this endpoint is accessed)
			msgString, err := dbactions.GetStringFromBucket(db, "pingResponses", "ping_second")
			if err != nil {
				log.Println(err.Error()) // FYI how do real error handling in handlerfunc?
			}

			log.Printf("retrieved message: %s\n", msgString)

			c.JSON(200, gin.H{ // parse this into JSON (and return?)
				"message": string(msgString),
			})
		},
	)

	// listen and serve on <instance localhost/IP>:<port>
	// accessible through the wm's eth0 IP in WSL
	err = r.Run(":" + apiPort)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
