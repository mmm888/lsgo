package main

import (
	"os"
	"strings"
)

type NormalFormat struct {
	Files string

	fList   []string
	options *Options
}

func NewNormalFormat(o *Options) *NormalFormat {
	return &NormalFormat{
		fList:   make([]string, 0, 5),
		options: o,
	}
}

func (n *NormalFormat) Execute(fInfos []os.FileInfo) {
	for _, fInfo := range fInfos {

		// Check hidden file
		var name string
		name = string(fInfo.Name()[0])
		if name == "." {
			continue
		}

		var info string
		info = getFileName(fInfo)
		n.fList = append(n.fList, info)
	}
	n.Files = strings.Join(n.fList, n.options.delimiter)
}
