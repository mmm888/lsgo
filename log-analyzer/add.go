/*
db を別パッケージにして、interface として利用
directory 以下のログファイルを全て add
long_l/server300/app.log > eclipse...
*/

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
)

type logFormat struct {
	Time     time.Time
	LogLevel string
	Host     string
	Actor    string
	Message  string
	CPU      float64
}

func (l *logFormat) GetAllText() string {
	return fmt.Sprintf("%s %s [%s] [%s] %s", l.Time.Format(logTimeFormat), l.LogLevel, l.Host, l.Actor, l.Message)
}

// parse `resource metric {CPU: 0.090}` > 0.090
func (l *logFormat) parseCPUUsage() error {

	var err error
	HeadMessage := "resource metric"
	TrimSymbol := " {}"
	CPUMessage := "CPU: "

	m := l.Message

	if !strings.HasPrefix(m, HeadMessage) {
		l.CPU = 0
		return nil
	}

	m = strings.TrimPrefix(m, HeadMessage)
	m = strings.Trim(m, TrimSymbol)
	m = strings.TrimPrefix(m, CPUMessage)

	l.CPU, err = strconv.ParseFloat(m, 64)
	if err != nil {
		return err
	}

	return nil
}

func (l *logFormat) insertDB(stmt *sql.Stmt) error {

	// at, loglevel, host, cpu usage, all text
	_, err := stmt.Exec(l.Time.Format(showTimeFormat), l.LogLevel, l.Host, l.CPU, l.GetAllText())
	if err != nil {
		return err
	}

	return nil
}

// parse `2018-04-01T00:00:00.094+0900 INFO [server107] [mesos-resource-actor] resource metric {CPU: 0.090}`
// > 2018-04-01T00:00:00.094+0900, INFO, server107, mesos-resource-actor, resource metric {CPU: 0.090}
func parseLine(line string) (*logFormat, error) {

	var l *logFormat
	var err error

	seps := strings.SplitN(line, " ", 5)

	for i := range seps {
		seps[i] = strings.Trim(seps[i], "[]")
	}

	t, err := time.Parse(logTimeFormat, seps[0])
	if err != nil {
		return nil, err
	}

	l = &logFormat{
		Time:     t,
		LogLevel: seps[1],
		Host:     seps[2],
		Actor:    seps[3],
		Message:  seps[4],
	}

	return l, nil
}

// if err1 is not nil, return err1. else, return err2.
func setErr(err1, err2 error) error {
	if err1 != nil {
		return err1
	}

	return err2
}

func fromFiletoDB(db *sql.DB, dbName, fp string) error {
	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("insert into %s (at, loglevel, host, cpu, alltext) values (?,?,?,?,?)", dbName)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	scanner := bufio.NewScanner(f)
	var errHandler error
	for scanner.Scan() {
		var err error

		l, err := parseLine(scanner.Text())
		errHandler = setErr(errHandler, err)

		err = l.parseCPUUsage()
		errHandler = setErr(errHandler, err)

		err = l.insertDB(stmt)
		errHandler = setErr(errHandler, err)
	}

	// roll back
	if errHandler != nil {
		tx.Rollback()
		return err
	}

	// roll back
	if err := scanner.Err(); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	tx.Commit()
	return nil
}

func Add(c *cli.Context) error {
	var err error

	logFile := c.String("f")
	dbPath := c.String("d")
	dbName := getFileNameWithoutExt(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	err = fromFiletoDB(db, dbName, logFile)
	if err != nil {
		return err
	}

	return nil
}
