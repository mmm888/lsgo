package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func getCheckCPU(db *sql.DB, cpu float64) error {

	query := fmt.Sprintf("select at, host from %s where cpu > %f", dbName, cpu)
	fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var t time.Time
		var host string
		err = rows.Scan(&t, &host)
		if err != nil {
			return err
		}
		fmt.Println(t.Format("2006-01-02 03:04:05"), host)
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func CPUCheck(f, d string, c int) error {

	logFile = f
	dbPath = d
	dbName = getFileNameWithoutExt(d)

	cpu := float64(c) / 100

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	err = getCheckCPU(db, cpu)
	if err != nil {
		return err
	}

	log.Print("Complete")
	return nil
}
