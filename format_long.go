package main

import "os"

type LongFormat struct {
	Total int
	List  []string

	options *Options
}

func NewLongFormat(o *Options) *LongFormat {
	return &LongFormat{options: o}
}

// Calculate Total, List
func (l *LongFormat) Execute(fInfos []os.FileInfo) {
	for _, fInfo := range fInfos {

		// Check hidden file
		var name string
		name = string(fInfo.Name()[0])
		if name == "." {
			continue
		}

		fi := NewFileInfo(fInfo)

		var info string
		info = fi.LongFormat(l.options)
		l.List = append(l.List, info)

		var size int
		size = fi.GetUsedBlockSize()
		l.Total += size
	}
}
