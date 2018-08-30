package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/urfave/cli"
)

type checkCPU struct {
	Time time.Time
	Host string
}

func getCheckCPU(db *sql.DB, dbName string, cpu float64) ([]checkCPU, error) {

	cs := make([]checkCPU, 0, 100)

	query := fmt.Sprintf("select at, host from %s where cpu > %f", dbName, cpu)
	rows, err := db.Query(query)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	for rows.Next() {
		var c checkCPU
		err = rows.Scan(&c.Time, &c.Host)
		if err != nil {
			return nil, err
		}

		cs = append(cs, c)
		//fmt.Println(c.Time.Format(showTimeFormat), c.Host)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func CPUCheck(c *cli.Context) error {

	dbPath := c.String("d")
	dbName := getFileNameWithoutExt(dbPath)
	cpu := float64(c.Int("c")) / 100

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	cs, err := getCheckCPU(db, dbName, cpu)
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
