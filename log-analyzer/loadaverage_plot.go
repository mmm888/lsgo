package main

import (
	"database/sql"
	"time"
)

type plotLoadAverages struct {
	start time.Time
	// 中央値を計算する最後の時間
	end     time.Time
	medians int
	avgs    []laData
	hosts   []string

	db     *sql.DB
	dbName string
}

func NewPlotLoadAverages(db *sql.DB, dbName string, medians int) (*plotLoadAverages, error) {

	s, e, hs, err := getCommonData(db, dbName, medians)
	if err != nil {
		return nil, err
	}

	// start, end は近似値を返す
	return &plotLoadAverages{
		start:   s.Round(time.Minute),
		end:     e.Round(time.Minute),
		avgs:    make([]laData, 0, 1000),
		medians: medians,
		hosts:   hs,
		db:      db,
		dbName:  dbName,
	}, nil

	return nil, nil
}

func (ps *plotLoadAverages) GetData() error {
	return nil
}

func (ps *plotLoadAverages) Output() error {
	return nil
}
