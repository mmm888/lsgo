package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
)

type showLoadAvarages struct {
	start time.Time
	// 中央値を計算する最後の時間
	end     time.Time
	medians int
	avgs    []showLoadAvarage
}

func NewShowLoadAvarages(db *sql.DB, dbName string, medians int) (*showLoadAvarages, error) {

	var err error

	row1 := db.QueryRow("select at from test order by at asc limit 1")
	var start time.Time
	err = row1.Scan(&start)
	if err != nil {
		return nil, err
	}

	row2 := db.QueryRow("select at from test order by at desc limit 1")
	var end time.Time
	err = row2.Scan(&end)
	if err != nil {
		return nil, err
	}

	// end - median(min)
	end = end.Add(time.Duration(-medians) * time.Minute)

	// start, end は近似値を返す
	return &showLoadAvarages{
		start:   start.Round(time.Minute),
		end:     end.Round(time.Minute),
		avgs:    make([]showLoadAvarage, 0, 1000),
		medians: medians,
	}, nil
}

func (ss *showLoadAvarages) getLoadAverage(db *sql.DB, dbName string) error {

	var err error
	var query string

	t := ss.start

	for !t.After(ss.end) {
		var count int
		var sum float64

		query = fmt.Sprintf("SELECT count(1) FROM test WHERE at > datetime('%s') AND at < datetime('%s', '+%d minutes') AND host = '%s' AND cpu != 0",
			t.Format(showTimeFormat), t.Format(showTimeFormat), ss.medians, "server101")
		row1 := db.QueryRow(query)
		err = row1.Scan(&count)
		if err != nil {
			return err
		}

		query = fmt.Sprintf("SELECT sum(cpu) FROM test WHERE at > datetime('%s') AND at < datetime('%s', '+%d minutes') AND host = '%s' AND cpu != 0",
			t.Format(showTimeFormat), t.Format(showTimeFormat), ss.medians, "server101")
		row2 := db.QueryRow(query)
		err = row2.Scan(&sum)
		if err != nil {
			return err
		}

		var s showLoadAvarage
		s = showLoadAvarage{
			Start:       t,
			End:         t.Add(time.Minute),
			LoadAverage: sum / float64(count),
			Host:        "server101",
		}

		ss.avgs = append(ss.avgs, s)

		t = t.Add(time.Minute)

		/*
			query := fmt.Sprintf("SELECT at, host, cpu FROM test WHERE at > datetime('%s') AND at < datetime('%s', '+%d minutes') AND host = '%s' AND cpu != 0",
				t.Format(showTimeFormat), t.Format(showTimeFormat), s.medians, "server101")
			rows, err := db.Query(query)
			if err != nil {
				return nil
			}
			defer rows.Close()

			for rows.Next() {
				var t time.Time
				var h string
				var c float64
				err = rows.Scan(&t, &h, &c)
				if err != nil {
					return err
				}

				fmt.Println(t, h, c)
			}

			err = rows.Err()
			if err != nil {
				return err
			}
		*/
	}

	return nil
}

type showLoadAvarage struct {
	Start       time.Time
	End         time.Time
	LoadAverage float64
	Host        string
}

func ShowLoadAverage(c *cli.Context) error {

	dbPath := c.String("d")
	dbName := getFileNameWithoutExt(dbPath)
	medians := c.Int("m")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	cs, err := NewShowLoadAvarages(db, dbName, medians)
	if err != nil {
		return err
	}

	err = cs.getLoadAverage(db, dbName)
	if err != nil {
		return err
	}

	avgs := cs.avgs
	if len(avgs) == 0 {
		fmt.Println("Nothing")
	} else {
		for i := range avgs {
			fmt.Printf("%s,%s,%.2f\n", avgs[i].Host, avgs[i].Start.Format(showTimeFormat), avgs[i].LoadAverage*100)
		}
	}

	return nil
}
