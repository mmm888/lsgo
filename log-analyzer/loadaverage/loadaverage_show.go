package loadaverage

import (
	"database/sql"
	"fmt"
)

type showLoadAvarages struct {
	avgs   map[string][]LAData
	median int
	db     *sql.DB
	table  string
}

func NewShowLoadAvarages(db *sql.DB, table string, median int) *showLoadAvarages {
	return &showLoadAvarages{
		avgs:   make(map[string][]LAData),
		median: median,
		db:     db,
		table:  table,
	}
}

func (ss *showLoadAvarages) GetData() error {
	var err error
	var query string

	query = fmt.Sprintf("select start, host, loadavg from %s where median = %d", ss.table, ss.median)
	rows, err := ss.db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var la LAData
		err = rows.Scan(&la.Start, &la.Host, &la.LoadAverage)
		if err != nil {
			return err
		}

		ss.avgs[la.Host] = append(ss.avgs[la.Host], la)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

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
