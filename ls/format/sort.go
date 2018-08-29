package format

import (
	"os"
	"sort"

	"github.com/mmm888/lsgo/option"
)

func sortList(o *option.Options, fi []os.FileInfo) {
	isTime := o.Time
	isReverse := o.Reverse

	if isTime {
		sort.Slice(fi, func(i, j int) bool {
			return fi[i].ModTime().After(fi[j].ModTime())
		})
	}

	if isReverse {
		for i, j := 0, len(fi)-1; i < j; i, j = i+1, j-1 {
			fi[i], fi[j] = fi[j], fi[i]
		}
	}
}
