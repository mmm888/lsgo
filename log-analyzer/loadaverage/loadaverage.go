package loadaverage

import (
	"database/sql"
	"time"
)

const (
	showTimeFormat = "2006-01-02 15:04:05"
)

type LoadAverage interface {
	GetData() error
	Output() error
}

type LAData struct {
	Start       time.Time
	End         time.Time
	LoadAverage float64
	Host        string
}

func (l *LAData) InsertDB(stmt *sql.Stmt, median int) error {
	// start, host, loadavg, median
	_, err := stmt.Exec(l.Start.Format(showTimeFormat), l.Host, l.LoadAverage, median)
	if err != nil {
		return err
	}

	return nil
}
