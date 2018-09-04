package cli

import (
	"database/sql"
	"fmt"

	"github.com/mmm888/mycmd/log-analyzer/loadaverage"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func LoadAvgMessage(c *cli.Context) error {
	dbPath := c.GlobalString("d")
	dbName := getFileNameWithoutExt(dbPath)
	logfileTable := fmt.Sprintf("%s_%s", dbName, logfileTableName)
	loadavgTable := fmt.Sprintf("%s_%s", dbName, loadavgTableName)
	median := c.GlobalInt("m")
	cpu := c.Int("c")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Wrap(err, "Error1: ")
	}
	defer db.Close()

	var s loadaverage.LoadAverage
	s = loadaverage.NewMessageLoadAverages(db, logfileTable, loadavgTable, median, cpu)

	err = s.GetData()
	if err != nil {
		return errors.Wrap(err, "Error2: ")
	}

	err = s.Output()
	if err != nil {
		return errors.Wrap(err, "Error3: ")
	}

	return nil
}
