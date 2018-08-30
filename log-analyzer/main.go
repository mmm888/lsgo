package main

import (
	"log"
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
				err := Add(c.String("f"), c.String("d"))
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "file, f",
					//Value: "app.log",
					Value: "log/log_s/server100/app.log",
				},
				cli.StringFlag{
					Name: "db, d",
					//Value: "test.db",
					Value: "test.db",
				},
			},
		},

		{
			Name:  "cpucheck",
			Usage: "Check threshold of CPU usage.",
			Action: func(c *cli.Context) error {
				err := CPUCheck(c.String("f"), c.String("d"), c.Int("c"))
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "file, f",
					//Value: "app.log",
					Value: "log/log_s/server100/app.log",
				},
				cli.StringFlag{
					Name: "db, d",
					//Value: "test.db",
					Value: "test.db",
				},
				cli.IntFlag{
					Name:  "cpu, c",
					Value: 95,
				},
			},
		},
	}

	app.Run(os.Args)
}
