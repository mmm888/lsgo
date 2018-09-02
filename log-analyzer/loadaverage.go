package main

import (
	"database/sql"
	"fmt"
	"time"
)

type LoadAverage interface {
	GetData() error
	Output() error
}

func getCommonData(db *sql.DB, dbName string, medians int) (start, end time.Time, hosts []string, err error) {
	var query string

	// Get Start time
	query = fmt.Sprintf("select at from %s order by at asc limit 1", dbName)
	row1 := db.QueryRow(query)
	err = row1.Scan(&start)
	if err != nil {
		return time.Time{}, time.Time{}, nil, err
	}

	// Get End time
	query = fmt.Sprintf("select at from %s order by at desc limit 1", dbName)
	row2 := db.QueryRow(query)
	err = row2.Scan(&end)
	if err != nil {
		return time.Time{}, time.Time{}, nil, err
	}
	// end - median(min)
	end = end.Add(time.Duration(-medians) * time.Minute)

	// Get hosts
	query = fmt.Sprintf("select host from %s group by host", dbName)
	row3, err := db.Query(query)
	if err != nil {
		return time.Time{}, time.Time{}, nil, err
	}
	defer row3.Close()

	for row3.Next() {
		var host string
		err = row3.Scan(&host)
		if err != nil {
			return time.Time{}, time.Time{}, nil, err
		}

		hosts = append(hosts, host)
	}
	err = row3.Err()
	if err != nil {
		return time.Time{}, time.Time{}, nil, err
	}

	return start, end, hosts, nil
}

func getLoadAverageFromDB(start, end time.Time, db *sql.DB, dbName string, medians int, host string) ([]laData, error) {
	var err error
	var query string

	avgs := make([]laData, 0, 100)
	t := start

	// end まで
	for !t.After(end) {
		var count int
		var sum float64

		query = fmt.Sprintf("SELECT count(1) FROM %s WHERE at > datetime('%s') AND at < datetime('%s', '+%d minutes') AND host = '%s' AND cpu != 0",
			dbName, t.Format(showTimeFormat), t.Format(showTimeFormat), medians, host)
		row1 := db.QueryRow(query)
		err = row1.Scan(&count)
		if err != nil {
			return nil, err
		}

		query = fmt.Sprintf("SELECT sum(cpu) FROM %s WHERE at > datetime('%s') AND at < datetime('%s', '+%d minutes') AND host = '%s' AND cpu != 0",
			dbName, t.Format(showTimeFormat), t.Format(showTimeFormat), medians, host)
		row2 := db.QueryRow(query)
		err = row2.Scan(&sum)
		if err != nil {
			return nil, err
		}

		s := laData{
			Start:       t,
			End:         t.Add(time.Minute),
			LoadAverage: sum / float64(count),
			Host:        host,
		}

		avgs = append(avgs, s)
		t = t.Add(time.Minute)
	}

	return avgs, nil
}
