package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
)

type checkUsage struct {
	Time time.Time
	Host string
}

func getCPUUsage(db *sql.DB, dbName string, cpu float64) ([]checkUsage, error) {

	cs := make([]checkUsage, 0, 100)

	query := fmt.Sprintf("select at, host from %s where cpu > %f", dbName, cpu)
	rows, err := db.Query(query)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	for rows.Next() {
		var c checkUsage
		err = rows.Scan(&c.Time, &c.Host)
		if err != nil {
			return nil, err
		}

		cs = append(cs, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func CPUUsage(c *cli.Context) error {

	dbPath := c.String("d")
	dbName := getFileNameWithoutExt(dbPath)
	cpu := float64(c.Int("c")) / 100

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	cs, err := getCPUUsage(db, dbName, cpu)
	if err != nil {
		return err
	}

	if len(cs) == 0 {
		fmt.Println("Nothing")
	} else {
		for i := range cs {
			fmt.Println(cs[i].Time.Format(showTimeFormat), cs[i].Host)
		}
	}

	return nil
}
