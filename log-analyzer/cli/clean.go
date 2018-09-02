package cli

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func Clean(c *cli.Context) error {
	var err error

	dbPath := c.String("d")
	dbName := getFileNameWithoutExt(dbPath)
	logfileTable := fmt.Sprintf("%s_%s", dbName, logfileTableName)
	loadavgTable := fmt.Sprintf("%s_%s", dbName, loadavgTableName)

	os.Remove(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Wrap(err, "Error1: ")
	}

	/*
		// primary key の設定
			sqlStmt := `
			create table foo (id integer not null primary key, name text);
			delete from foo;
			`
	*/

	var sqlStmt string

	sqlStmt = fmt.Sprintf(`
	create table %s (at datetime, loglevel text, host text, cpu float, alltext text);
	delete from %s;
	`, logfileTable, logfileTable)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return errors.Wrap(err, "Error2: ")
	}

	sqlStmt = fmt.Sprintf(`
	create table %s (start datetime, host text, loadavg float, median int);
	delete from %s;
	`, loadavgTable, loadavgTable)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return errors.Wrap(err, "Error3: ")
	}

	return nil
}
