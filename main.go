package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mmm888/lsgo/format"
	"github.com/mmm888/lsgo/option"
)

func getFileInfos(o *option.Options) ([]os.FileInfo, error) {
	root := o.Root
	isAll := o.All

	f, err := os.Stat(root)
	// Error: No Such file or directory
	if err != nil {
		msg := fmt.Sprintf("%s: No such file or directory", root)
		return nil, errors.New(msg)
	}

	var list []os.FileInfo
	if f.IsDir() {
		list = make([]os.FileInfo, 0, 10)

		// Check current directory, parent directory
		if isAll {
			cd, _ := os.Stat(".")
			pd, _ := os.Stat("..")
			list = append(list, cd, pd)
		}

		fi, err := ioutil.ReadDir(root)
		if err != nil {
			return nil, err
		}

		if isAll {
			list = append(list, fi...)
		} else {
			for _, i := range fi {

				// Check hidden file
				name := string(i.Name()[0])
				if name == "." {
					continue
				}

				list = append(list, i)
			}
		}

	} else {
		list = append(list, f)
	}

	return list, nil
}

func main() {

	var options *option.Options
	options = option.CreateOptions()

	fInfos, err := getFileInfos(options)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	var frt format.Format
	if options.Long {
		frt = format.NewLongFormat(options)
	} else {
		frt = format.NewNormalFormat(options)
	}

	frt.Execute(fInfos)
}
