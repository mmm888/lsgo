package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

type laData struct {
	Start       time.Time
	End         time.Time
	LoadAverage float64
	Host        string
}

type showLoadAvarages struct {
	medians int
	avgs    map[string][]laData

	db     *sql.DB
	dbName string
}

func NewShowLoadAvarages(db *sql.DB, dbName string, medians int) (*showLoadAvarages, error) {

	s, e, hs, err := getCommonData(db, dbName, medians)
	if err != nil {
		return nil, err
	}

	// 1分で近似値を返す
	s = s.Round(time.Minute)
	e = e.Round(time.Minute)

	var wg sync.WaitGroup
	avgs := make(map[string][]laData)
	for i := range hs {
		host := hs[i]
		wg.Add(1)
		go func() {
			a, err := getLoadAverageFromDB(s, e, db, dbName, medians, host)
			defer wg.Done()
			if err != nil {
				log.Print(err)
			}
			avgs[host] = a
		}()
	}
	wg.Wait()

	return &showLoadAvarages{
		avgs:    avgs,
		medians: medians,

		db:     db,
		dbName: dbName,
	}, nil
}

func (ss *showLoadAvarages) GetData() error {

	/*
		var err error
		var query string

		// すべての host について
		for i := range ss.hosts {
			host := ss.hosts[i]

			getLoadAverageFromDB(ss, host)

		}
	*/

	return nil
}

func (ss *showLoadAvarages) Output() error {

	avgs := ss.avgs
	if len(avgs) == 0 {
		fmt.Println("Nothing")
	} else {
		for host, avgs := range avgs {
			for i := range avgs {
				la := avgs[i]
				fmt.Printf("%s,%s,%.2f\n", host, la.Start.Format(showTimeFormat), la.LoadAverage*100)
			}
		}
	}

	return nil
}
