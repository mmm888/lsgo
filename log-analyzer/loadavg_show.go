package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func LoadAvgShow(c *cli.Context) error {

	dbPath := c.GlobalString("d")
	dbName := getFileNameWithoutExt(dbPath)
	medians := c.GlobalInt("m")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Wrap(err, "Error1: ")
	}
	defer db.Close()

	var s LoadAverage
	s, err = NewShowLoadAvarages(db, dbName, medians)
	if err != nil {
		return errors.Wrap(err, "Error2: ")
	}

	err = s.GetData()
	if err != nil {
		return errors.Wrap(err, "Error3: ")
	}

	err = s.Output()
	if err != nil {
		return errors.Wrap(err, "Error4: ")
	}

	return nil
}
