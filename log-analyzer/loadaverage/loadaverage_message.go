package loadaverage

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type messageLoadAvarages struct {
	median       int
	messages     map[string][]string
	cpuusage     int
	db           *sql.DB
	logfileTable string
	loadavgTable string
}

func NewMessageLoadAverages(db *sql.DB, logfile, loadavg string, median, cpu int) *messageLoadAvarages {
	return &messageLoadAvarages{
		messages:     make(map[string][]string),
		median:       median,
		cpuusage:     cpu,
		db:           db,
		logfileTable: logfile,
		loadavgTable: loadavg,
	}
}

func (ms *messageLoadAvarages) GetData() error {
	var err error
	var query string

	query = fmt.Sprintf("select start, host from %s where median = %d AND loadavg > %.3f",
		ms.loadavgTable, ms.median, float64(ms.cpuusage)/100)
	log.Println(query)
	rows, err := ms.db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	avgs := make(map[string][]LAData)
	for rows.Next() {
		var la LAData
		err = rows.Scan(&la.Start, &la.Host)
		if err != nil {
			return err
		}
		la.End = la.Start.Add(time.Duration(ms.median))

		avgs[la.Host] = append(avgs[la.Host], la)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	for host := range avgs {
		for _, data := range avgs[host] {
			query = fmt.Sprintf("select alltext from %s where host = '%s' AND at > datetime('%s') AND at < datetime('%s', '+%d minutes') AND (loglevel = 'WARNING' OR loglevel = 'ERROR')",
				ms.logfileTable, host, data.Start.Format(showTimeFormat), data.End.Format(showTimeFormat), ms.median)
			//log.Print(query)
			rows, err := ms.db.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var msg string
				err = rows.Scan(&msg)
				if err != nil {
					return err
				}

				ms.messages[host] = append(ms.messages[host], msg)
			}
		}
	}

	return nil
}

func (ms *messageLoadAvarages) Output() error {

	arr := []string{"a", "b", "c", "a"}
	m := make(map[string]bool)
	uniq := []string{}

	for _, ele := range arr {
		if !m[ele] {
			m[ele] = true
			uniq = append(uniq, ele)
		}
	}

	messagesList := ms.messages
	if len(messagesList) == 0 {
		fmt.Println("Nothing")
	} else {
		for host, messages := range messagesList {
			log.Printf("%s's log", host)
			m := make(map[string]bool)
			//uniq := make([]string, 0, 100)

			for _, message := range messages {
				if !m[message] {
					m[message] = true
					// ホントはログ本文のみ
					fmt.Println(message)
				}
			}
		}
	}

	return nil
}
