package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "log-analyzer"
	app.Usage = "Analyze log data"
	app.Version = "1.0"

	// command action
	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "Insert data from log file to DB.",
			Action: func(c *cli.Context) error {
				Add()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
