package cli

import "path/filepath"

const (
	logTimeFormat  = "2006-01-02T15:04:05Z0700"
	showTimeFormat = "2006-01-02 15:04:05"

	logfileTableName = "logfile"
	loadavgTableName = "loadavg"
)

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

// if err1 is not nil, return err1. else, return err2.
func setErr(err1, err2 error) error {
	if err1 != nil {
		return err1
	}

	return err2
}
