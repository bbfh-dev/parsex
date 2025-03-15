package parsex

var HelpFlag = FlagOption{
	Name:     "help",
	Keywords: []string{"help", "h"},
	Desc:     "Print help and exit",
}

var VersionFlag = FlagOption{
	Name:     "version",
	Keywords: []string{"version", "V"},
	Desc:     "Print version and exit",
}

type FlagOption struct {
	Name     string
	Keywords []string
	Desc     string
}

func (opt FlagOption) Id() string {
	return opt.Name
}

func (opt FlagOption) Match() []string {
	return opt.Keywords
}

func (opt FlagOption) Describe() string {
	return opt.Desc
}
