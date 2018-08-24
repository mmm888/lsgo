package format

import (
	"os"
	"text/template"

	"github.com/mmm888/lsgo/fileinfo"
	"github.com/mmm888/lsgo/option"
)

const (
	templateLongFormat = `total {{ .Total }} 
{{ range .List -}}
{{ . }}
{{ end -}}
`
)

type LongFormat struct {
	Total int
	List  []string

	options *option.Options
	tmpl    *template.Template
}

func NewLongFormat(o *option.Options) *LongFormat {
	return &LongFormat{
		tmpl:    template.Must(template.New("long").Parse(templateLongFormat)),
		options: o}
}

// template Execute
func (l *LongFormat) Execute(fInfos []os.FileInfo) {
	for _, fInfo := range fInfos {

		fi := fileinfo.NewFileInfo(fInfo)

		var info string
		info = fi.LongFormat(l.options)
		l.List = append(l.List, info)

		var size int
		size = fi.GetUsedBlockSize()
		l.Total += size
	}

	l.tmpl.Execute(os.Stdout, l)
}
