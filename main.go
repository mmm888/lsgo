package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	templateNormal = `
{{- . }}
`

	templateLongFormat = `total {{ .Total }} 
{{ range .List -}}
{{ . }}
{{ end -}}
`
)

var (
	delimiter = " "

	tmpl *template.Template
)

type LF struct {
	Total int
	List  []string
}

func NewLF() *LF {
	return &LF{}
}

func main() {
	options := NewOption()
	options.Init()
	options.Check()

	fInfos, err := ioutil.ReadDir(options.root)
	if err != nil {
		log.Print("error1")
	}

	if options.long {

		lf := NewLF()
		tmpl = template.Must(template.New("long").Parse(templateLongFormat))
		for _, fInfo := range fInfos {

			// Check hidden file
			var name string
			name = string(fInfo.Name()[0])
			if name == "." {
				continue
			}

			fi := NewFileInfo(fInfo)

			var info string
			info = fi.LongFormat(options)
			lf.List = append(lf.List, info)

			var size int
			size = fi.GetUsedBlockSize()
			lf.Total += size
		}
		tmpl.Execute(os.Stdout, lf)

	} else {
		fList := make([]string, 0, 5)
		tmpl = template.Must(template.New("normal").Parse(templateNormal))
		for _, fInfo := range fInfos {

			// Check hidden file
			var name string
			name = string(fInfo.Name()[0])
			if name == "." {
				continue
			}

			var info string
			info = getFileName(fInfo)
			fList = append(fList, info)
		}
		tmpl.Execute(os.Stdout, strings.Join(fList, delimiter))

	}
}
