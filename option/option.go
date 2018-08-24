package option

import "flag"

type Options struct {
	Long  bool
	All   bool
	Human bool
	Root  string

	// Form normal format
	Delimiter string

	// TODO: rm
	Args []string
}

func CreateOptions() *Options {
	o := NewOptions()
	o.Init()
	o.Check()

	return o
}

func NewOptions() *Options {
	return &Options{
		Root:      ".",
		Delimiter: " ",
	}
}

func (o *Options) Init() {
	flag.BoolVar(&o.All, "a", false, "Include directory entries whose names begin with a dot (.).")
	flag.BoolVar(&o.Human, "h", false, "File size in human readable format.")
	flag.BoolVar(&o.Long, "l", false, "List in long format.")
	flag.Parse()

	o.Args = flag.Args()
}

func (o *Options) Check() {
	if len(o.Args) > 0 {
		o.Root = o.Args[0]
	}
}
