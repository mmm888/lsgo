package cli

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mmm888/mycmd/log-analyzer/loadaverage"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type plotData struct {
	Host        string
	Time        time.Time
	LoadAverage float64
}

func LoadAvgPlot(c *cli.Context) error {

	dbPath := c.GlobalString("d")
	dbName := getFileNameWithoutExt(dbPath)
	tableName := fmt.Sprintf("%s_%s", dbName, loadavgTableName)
	median := c.GlobalInt("m")
	output := c.String("o")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Wrap(err, "Error1: ")
	}
	defer db.Close()

	var s loadaverage.LoadAverage
	s = loadaverage.NewPlotLoadAverages(db, tableName, median, output)

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
