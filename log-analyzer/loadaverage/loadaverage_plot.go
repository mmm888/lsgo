package loadaverage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	chart "github.com/wcharczuk/go-chart"
)

type plotLoadAverages struct {
	avgs   map[string][]LAData
	median int
	db     *sql.DB
	table  string
	output string
}

func NewPlotLoadAverages(db *sql.DB, table string, median int, output string) *plotLoadAverages {
	return &plotLoadAverages{
		avgs:   make(map[string][]LAData),
		median: median,
		db:     db,
		table:  table,
		output: output,
	}
}

func (ps *plotLoadAverages) GetData() error {
	var err error
	var query string

	query = fmt.Sprintf("select start, host, loadavg from %s where median = %d", ps.table, ps.median)
	rows, err := ps.db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var la LAData
		err = rows.Scan(&la.Start, &la.Host, &la.LoadAverage)
		if err != nil {
			return err
		}
		la.End = la.Start.Add(time.Duration(ps.median))

		ps.avgs[la.Host] = append(ps.avgs[la.Host], la)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func valueFormatter(v interface{}) string {
	dateFormat := "01-02 15:04"

	if typed, isTyped := v.(time.Time); isTyped {
		return typed.Format(dateFormat)
	}
	if typed, isTyped := v.(int64); isTyped {
		return time.Unix(0, typed).Format(dateFormat)
	}
	if typed, isTyped := v.(float64); isTyped {
		return time.Unix(0, int64(typed)).Format(dateFormat)
	}
	return ""
}

func (ps *plotLoadAverages) Output() error {

	numSeries := len(ps.avgs)
	series := make([]chart.Series, numSeries)

	var i int
	for host, laDatas := range ps.avgs {
		numValues := len(laDatas)
		xValues := make([]time.Time, numValues)
		yValues := make([]float64, numValues)

		for i, LAData := range laDatas {
			xValues[i] = LAData.Start
			yValues[i] = LAData.LoadAverage
		}

		series[i] = chart.TimeSeries{
			Name:    fmt.Sprintf("%s.value", host),
			XValues: xValues,
			YValues: yValues,
		}
		i = i + 1
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:           "Time",
			NameStyle:      chart.StyleShow(),
			Style:          chart.StyleShow(),
			ValueFormatter: valueFormatter,
		},
		YAxis: chart.YAxis{
			Name:      "Load Average",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		Series: series,
	}

	/*
		graph.Elements = []chart.Renderable{
			chart.Legend(&graph),
		}
	*/

	f, err := os.Create(ps.output)
	if err != nil {
		return err
	}
	defer f.Close()

	graph.Render(chart.PNG, f)

	return err
}
