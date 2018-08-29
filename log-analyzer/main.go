package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath  = "./test.db"
	dbname  = "test"
	logFile = "./log_s/server100/app.log"

	timeFormat = "2006-01-02T15:04:05Z0700"
	splitDel   = " "
	splitNum   = 5
)

type logFormat struct {
	Time     time.Time
	LogLevel string
	Host     string
	Actor    string
	Message  string
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

func insertDB(stmt *sql.Stmt, l *logFormat) error {

	_, err := stmt.Exec(l.Time.Format("2006-01-02 03:04:05"), l.LogLevel, l.Host, l.Actor, l.Message)
	if err != nil {
		return err
	}

	return nil
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
	stmt, err := tx.Prepare("insert into test (at,loglevel,host,actor,message) values (?,?,?,?,?)")
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

		err = insertDB(stmt, l)
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
	create table %s (at datetime, loglevel text, host text, actor text, message text);
	delete from %s;
	`, dbname, dbname)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil, err
	}

	return db, nil
}

func main() {
	var err error

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = fromFiletoDB(db, logFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Complete")
}
