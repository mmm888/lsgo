package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
)

func Clean(c *cli.Context) error {
	var err error

	dbPath := c.String("d")
	dbName := getFileNameWithoutExt(dbPath)

	os.Remove(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
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
		return err
	}

	return nil
}
