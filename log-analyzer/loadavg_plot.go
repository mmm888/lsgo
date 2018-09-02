package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type plotData struct {
	Host        string
	Time        time.Time
	LoadAverage float64
}

/*
func fromFiletoPlotData(fp string) ([]plotData, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pds := make([]plotData, 0, 100)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var p plotData

		s := strings.Split(scanner.Text(), ",")
		pf, _ := strconv.ParseFloat(s[2], 64)
		pt, _ := time.Parse(showTimeFormat, s[1])
		p = plotData{
			Host:        s[0],
			Time:        pt,
			LoadAverage: pf,
		}

		pds = append(pds, p)
	}

	return pds, nil
}
*/

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

	var s LoadAverage
	s = NewPlotLoadAverages(db, tableName, median, output)

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
