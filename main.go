package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	// bolt "go.etcd.io/bbolt"
)

var dbPort string
var apiPort string

func main() {

	app := &cli.App{
		Name:   "apiserver",
		Usage:  "Hi hello yes here should be information", // TODO
		Action: apiserv,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dbPort",
				Aliases:     []string{"db"},
				Usage:       "assign `PORT` to communicate with the database on",
				Destination: &dbPort,
				Value:       "4711",
			},
			&cli.StringFlag{
				Name:        "apiPort",
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

func apiserv(c *cli.Context) error {
	log.Printf("Server started")
	log.Printf("\tCommunicating with db on port: %s\n", dbPort)
	log.Printf("\tAPI accessible on port: %s\n", apiPort)
	fmt.Println()

	r := gin.Default()
	r.GET( // GET request on /ping endpoint
		"/ping", // (relative) PATH
		func(c *gin.Context) { // handlerfunc (what to do when this endpoint is accessed)
			c.JSON(200, gin.H{ // parse this into JSON (and return?)
				"message": "you rang?",
			})
		},
	)
	err := r.Run(":" + apiPort) // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
