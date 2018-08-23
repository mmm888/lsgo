package main

import "flag"

type Options struct {
	args  []string
	long  bool
	all   bool
	human bool
	root  string
}

func NewOption() *Options {
	return &Options{}
}

func (o *Options) Init() {
	flag.BoolVar(&o.all, "a", false, "Include directory entries whose names begin with a dot (.).")
	flag.BoolVar(&o.human, "h", false, "File size in human readable format.")
	flag.BoolVar(&o.long, "l", false, "List in long format.")
	flag.Parse()

	o.args = flag.Args()
}

func (o *Options) Check() {
	if len(o.args) < 1 {
		o.root = "."
	} else {
		o.root = o.args[0]
	}
}
