package cli

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mmm888/mycmd/log-analyzer/loadaverage"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func getCommonData(db *sql.DB, table string, median int) (start, end time.Time, hosts []string, err error) {
	var query string

	// Get Start time
	query = fmt.Sprintf("select at from %s order by at asc limit 1", table)
	fmt.Printf("Exec sql: %s\n", query)
	row1 := db.QueryRow(query)
	err = row1.Scan(&start)
	if err != nil {
		return time.Time{}, time.Time{}, nil, err
	}

	// Get End time
	query = fmt.Sprintf("select at from %s order by at desc limit 1", table)
	fmt.Printf("Exec sql: %s\n", query)
	row2 := db.QueryRow(query)
	err = row2.Scan(&end)
	if err != nil {
		return time.Time{}, time.Time{}, nil, err
	}
	// end - median(min)
	end = end.Add(time.Duration(-median) * time.Minute)

	// Get hosts
	query = fmt.Sprintf("select host from %s group by host", table)
	fmt.Printf("Exec sql: %s\n", query)
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

func getLoadAverageFromDB(start, end time.Time, db *sql.DB, table string, median int, host string) ([]loadaverage.LAData, error) {
	var err error
	var query string

	avgs := make([]loadaverage.LAData, 0, 100)
	t := start

	// end まで
	for !t.After(end) {
		var count int
		var sum float64

		query = fmt.Sprintf("SELECT count(1) FROM %s WHERE at > datetime('%s') AND at < datetime('%s', '+%d minutes') AND host = '%s' AND cpu != 0",
			table, t.Format(showTimeFormat), t.Format(showTimeFormat), median, host)
		fmt.Printf("Exec sql: %s\n", query)
		row1 := db.QueryRow(query)
		err = row1.Scan(&count)
		if err != nil {
			return nil, err
		}

		query = fmt.Sprintf("SELECT sum(cpu) FROM %s WHERE at > datetime('%s') AND at < datetime('%s', '+%d minutes') AND host = '%s' AND cpu != 0",
			table, t.Format(showTimeFormat), t.Format(showTimeFormat), median, host)
		fmt.Printf("Exec sql: %s\n", query)
		row2 := db.QueryRow(query)
		err = row2.Scan(&sum)
		if err != nil {
			return nil, err
		}

		s := loadaverage.LAData{
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

func fromDataToLoadavg(db *sql.DB, exportTable, importTable string, median int) error {

	s, e, hs, err := getCommonData(db, exportTable, median)
	if err != nil {
		return err
	}

	// 1分で近似値を返す
	s = s.Round(time.Minute)
	e = e.Round(time.Minute)

	var wg sync.WaitGroup
	mutex := &sync.Mutex{}
	avgs := make(map[string][]loadaverage.LAData)
	for i := range hs {
		host := hs[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			a, err := getLoadAverageFromDB(s, e, db, exportTable, median, host)
			if err != nil {
				log.Print(err)
			}

			mutex.Lock()
			avgs[host] = a
			mutex.Unlock()
		}()
	}
	wg.Wait()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("insert into %s (start, host, loadavg, median) values (?,?,?,?)", importTable)
	fmt.Printf("Exec sql: %s\n", query)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var errHandler error
	for host := range avgs {
		for _, la := range avgs[host] {
			var err error

			err = la.InsertDB(stmt, median)
			errHandler = setErr(errHandler, err)
		}
	}

	// roll back
	if errHandler != nil {
		tx.Rollback()
		return err
	}

	// commit
	tx.Commit()

	return nil
}

func AddLoadAvg(c *cli.Context) error {
	var err error

	dbPath := c.GlobalString("d")
	dbName := getFileNameWithoutExt(dbPath)
	logfileTable := fmt.Sprintf("%s_%s", dbName, logfileTableName)
	loadavgTable := fmt.Sprintf("%s_%s", dbName, loadavgTableName)
	median := c.Int("m")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Wrap(err, "Error1: ")
	}
	defer db.Close()

	err = fromDataToLoadavg(db, logfileTable, loadavgTable, median)
	if err != nil {
		return errors.Wrap(err, "Error2: ")
	}

	return nil
}
