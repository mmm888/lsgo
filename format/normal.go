package format

import (
	"os"
	"strings"
	"text/template"

	"github.com/mmm888/lsgo/fileinfo"
	"github.com/mmm888/lsgo/option"
)

const (
	templateNormal = `
{{- .Files }}
`
)

type NormalFormat struct {
	Files string

	fList   []string
	tmpl    *template.Template
	options *option.Options
}

func NewNormalFormat(o *option.Options) *NormalFormat {
	return &NormalFormat{
		fList:   make([]string, 0, 5),
		tmpl:    template.Must(template.New("normal").Parse(templateNormal)),
		options: o,
	}
}

// template Execute
func (n *NormalFormat) Execute(fInfos []os.FileInfo) {
	for _, fInfo := range fInfos {

		var info string
		info = fileinfo.GetFileName(fInfo)
		n.fList = append(n.fList, info)
	}
	n.Files = strings.Join(n.fList, n.options.Delimiter)

	n.tmpl.Execute(os.Stdout, n)
}
