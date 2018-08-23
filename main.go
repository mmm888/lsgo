package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

var (
	delimiter = " "
)

func main() {
	options := NewOption()
	options.Init()
	options.Check()

	fInfos, err := ioutil.ReadDir(options.root)
	if err != nil {
		log.Print("error1")
	}

	fList := make([]string, 0, 5)
	for _, fInfo := range fInfos {

		var info string

		// long format
		if options.long {
			// TODO: 改行コード
			delimiter = "\n"

			// total [total XX]
			info = longFormat(fInfo)
		} else {
			info = fInfo.Name()
		}

		fList = append(fList, info)
	}

	fmt.Println(strings.Join(fList, delimiter))
}
