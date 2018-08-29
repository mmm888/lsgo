package option

import "flag"

type Options struct {
	All     bool
	Human   bool
	Long    bool
	Reverse bool
	Time    bool
	Roots   []string

	// For normal format
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
		Roots:     []string{"."},
		Delimiter: " ",
	}
}

func (o *Options) Init() {
	flag.BoolVar(&o.All, "a", false, "Include directory entries whose names begin with a dot (.).")
	flag.BoolVar(&o.Human, "h", false, "File size in human readable format.")
	flag.BoolVar(&o.Long, "l", false, "List in long format.")
	flag.BoolVar(&o.Reverse, "r", false, "Reverse the order of the sort.")
	flag.BoolVar(&o.Time, "t", false, "Sort by time modified.")
	flag.Parse()

	o.Args = flag.Args()
}

func (o *Options) Check() {
	if len(o.Args) > 0 {
		o.Roots = o.Args
	}
}
