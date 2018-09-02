package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func LoadAvgShow(c *cli.Context) error {

	dbPath := c.GlobalString("d")
	dbName := getFileNameWithoutExt(dbPath)
	tableName := fmt.Sprintf("%s_%s", dbName, loadavgTableName)
	median := c.GlobalInt("m")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Wrap(err, "Error1: ")
	}
	defer db.Close()

	var s LoadAverage
	s = NewShowLoadAvarages(db, tableName, median)

	err = s.GetData()
	if err != nil {
		return errors.Wrap(err, "Error2: ")
	}

	err = s.Output()
	if err != nil {
		return errors.Wrap(err, "Error3: ")
	}

	return nil
}
