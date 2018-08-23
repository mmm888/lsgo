package main

import "flag"

type Options struct {
	args  []string
	long  bool
	all   bool
	human bool
	root  string

	// Form normal format
	delimiter string
}

func CreateOptions() *Options {
	o := NewOptions()
	o.Init()
	o.Check()

	return o
}

func NewOptions() *Options {
	return &Options{
		delimiter: " ",
	}
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
