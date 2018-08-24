package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mmm888/lsgo/format"
	"github.com/mmm888/lsgo/option"
)

func getFileInfos(o *option.Options) ([][]os.FileInfo, error) {
	roots := o.Roots
	isAll := o.All

	lists := make([][]os.FileInfo, 0, 5)
	for _, root := range roots {
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

		lists = append(lists, list)
	}

	return lists, nil
}

func main() {

	var err error

	var options *option.Options
	options = option.CreateOptions()

	fInfosList, err := getFileInfos(options)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	for _, fInfos := range fInfosList {
		var frt format.Format
		if options.Long {
			frt = format.NewLongFormat(options)
		} else {
			frt = format.NewNormalFormat(options)
		}

		err = frt.Execute(fInfos)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
}
