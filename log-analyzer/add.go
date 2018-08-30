/*
db を別パッケージにして、interface として利用
*/

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	timeFormat = "2006-01-02T15:04:05Z0700"
	splitDel   = " "
	splitNum   = 5
)

var (
	dbPath  string
	dbName  string
	logFile string
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
	return fmt.Sprintf("%s %s [%s] [%s] %s", l.Time.Format(timeFormat), l.LogLevel, l.Host, l.Actor, l.Message)
}

// parse `resource metric {CPU: 0.090}`
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
	_, err := stmt.Exec(l.Time.Format("2006-01-02 03:04:05"), l.LogLevel, l.Host, l.CPU, l.GetAllText())
	if err != nil {
		return err
	}

	return nil
}

func parseLine(line string) (*logFormat, error) {

	var l *logFormat
	var err error

	seps := strings.SplitN(line, splitDel, splitNum)

	for i := range seps {
		seps[i] = strings.Trim(seps[i], "[]")
	}

	t, err := time.Parse(timeFormat, seps[0])
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

func fromFiletoDB(db *sql.DB, fp string) error {
	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into test (at, loglevel, host, cpu, alltext) values (?,?,?,?,?)")
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

func initDB() (*sql.DB, error) {
	os.Remove(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	/*
		sqlStmt := `
		create table foo (id integer not null primary key, name text);
		delete from foo;
		`
	*/

	sqlStmt := fmt.Sprintf(`
	create table %s (at datetime, loglevel text, host text, cpu float, alltext text);
	delete from %s;
	`, dbName, dbName)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil, err
	}

	return db, nil
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func Add(f, d string) error {
	var err error

	logFile = f
	dbPath = d
	dbName = getFileNameWithoutExt(d)

	db, err := initDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = fromFiletoDB(db, logFile)
	if err != nil {
		return err
	}

	log.Print("Complete")
	return nil
}
