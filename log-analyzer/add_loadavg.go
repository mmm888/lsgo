package main

import (
	"database/sql"

	"github.com/urfave/cli"
)

func AddLoadAvg(c *cli.Context) error {
	var err error

	dbPath := c.GlobalString("d")
	//dbName := getFileNameWithoutExt(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	return nil
}
