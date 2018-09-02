package main

import (
	"fmt"
	"log"
	"os"

	mycli "github.com/mmm888/mycmd/log-analyzer/cli"
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
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Insert data to DB.",
			Subcommands: []cli.Command{
				{
					Name:  "logfile",
					Usage: "Insert data from log file.",
					Action: func(c *cli.Context) error {
						err := mycli.AddLogFile(c)
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
					},
				},
				{
					Name:  "loadavg",
					Usage: "Insert data of load average.",
					Action: func(c *cli.Context) error {
						err := mycli.AddLoadAvg(c)
						if err != nil {
							log.Fatal(err)
						}

						fmt.Println("Completed")
						return nil
					},
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "median, m",
							Value: 5,
						},
					},
				},
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
			Name:  "clean",
			Usage: "Reset DB data.",
			Action: func(c *cli.Context) error {
				err := mycli.Clean(c)
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
			Name:  "cpuusage",
			Usage: "Show threshold of CPU usage.",
			Action: func(c *cli.Context) error {
				err := mycli.CPUUsage(c)
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
			Name:    "loadavg",
			Aliases: []string{"l"},
			Usage:   "CPU load average",
			Subcommands: []cli.Command{
				{
					Name:  "show",
					Usage: "Show CPU load average.",
					Action: func(c *cli.Context) error {
						err := mycli.LoadAvgShow(c)
						if err != nil {
							log.Fatal(err)
						}

						return nil
					},
				},
				{
					Name:  "plot",
					Usage: "Plot CPU load average.",
					Action: func(c *cli.Context) error {
						err := mycli.LoadAvgPlot(c)
						if err != nil {
							log.Fatal(err)
						}

						return nil
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "output, o",
							Value: "test.png",
						},
					},
				},
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "db, d",
					//Value: "test.db",
					Value: "test.db",
				},
				cli.IntFlag{
					Name:  "median, m",
					Value: 5,
				},
			},
		},
	}

	app.Run(os.Args)
}
