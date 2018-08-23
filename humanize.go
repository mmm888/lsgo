package main

import (
	"fmt"
	"math"
)

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64) string {

	var base float64
	base = 1000
	sizes := []string{"B", "K", "M", "G", "T", "P", "E"}

	if s < 10 {
		return fmt.Sprintf("%dB", s)
	}

	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	//val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	val := float64(s) / math.Pow(base, e)
	f := "%.0f%s"
	if val < 10 {
		f = "%.1f%s"
	}

	return fmt.Sprintf(f, val, suffix)
}
