package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

const (
	templateNormal = `
{{- .Files }}
`

	templateLongFormat = `total {{ .Total }} 
{{ range .List -}}
{{ . }}
{{ end -}}
`
)

var (
	tmpl *template.Template
)

func main() {

	var options *Options
	options = CreateOptions()

	fInfos, err := ioutil.ReadDir(options.root)
	if err != nil {
		log.Print("error1")
	}

	if options.long {

		lf := NewLongFormat(options)
		lf.Execute(fInfos)
		tmpl = template.Must(template.New("long").Parse(templateLongFormat))
		tmpl.Execute(os.Stdout, lf)

	} else {

		nf := NewNormalFormat(options)
		nf.Execute(fInfos)
		tmpl = template.Must(template.New("normal").Parse(templateNormal))
		tmpl.Execute(os.Stdout, nf)

	}
}
