package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

const (
	logTimeFormat  = "2006-01-02T15:04:05Z0700"
	showTimeFormat = "2006-01-02 15:04:05"
)

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

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
				err := Add(c)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Completed")
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
			Name:  "clean",
			Usage: "Reset DB data.",
			Action: func(c *cli.Context) error {
				err := Clean(c)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Completed")
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "db, d",
					//Value: "test.db",
					Value: "test.db",
				},
			},
		},

		{
			Name:  "checkusage",
			Usage: "Check threshold of CPU usage.",
			Action: func(c *cli.Context) error {
				err := CheckUsage(c)
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
			Flags: []cli.Flag{
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

		{
			Name:  "showloadavg",
			Usage: "Show CPU load average.",
			Action: func(c *cli.Context) error {
				err := ShowLoadAverage(c)
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "db, d",
					//Value: "test.db",
					Value: "test.db",
				},
				cli.IntFlag{
					Name:  "medians, m",
					Value: 5,
				},
			},
		},

		{
			Name:  "plot",
			Usage: "Plot CPU load average.",
			Action: func(c *cli.Context) error {
				err := PlotLoadAverage(c)
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Value: "test.csv",
				},
			},
		},
	}

	app.Run(os.Args)
}
