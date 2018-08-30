package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
	chart "github.com/wcharczuk/go-chart"
)

type plotData struct {
	Host        string
	Time        time.Time
	LoadAverage float64
}

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

func plotLoadAvg(pds []plotData, fp string) error {

	numSeries := 1
	numValues := len(pds)

	series := make([]chart.Series, numSeries)

	for i := 0; i < numSeries; i++ {

		xValues := make([]time.Time, numValues)
		yValues := make([]float64, numValues)

		for j := 0; j < numValues; j++ {
			xValues[j] = pds[j].Time
			yValues[j] = pds[j].LoadAverage
		}

		series[i] = chart.TimeSeries{
			Name:    fmt.Sprintf("%s.value", pds[i].Host),
			XValues: xValues,
			YValues: yValues,
		}
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

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	graph.Render(chart.PNG, f)
	return err
}

func PlotLoadAverage(c *cli.Context) error {

	var err error

	dataFile := c.String("f")
	plotFile := "test.png"

	pds, err := fromFiletoPlotData(dataFile)
	if err != nil {
		return err
	}

	err = plotLoadAvg(pds, plotFile)
	if err != nil {
		return err
	}

	fmt.Println("Completed")

	return nil
}
