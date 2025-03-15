package parsex

type ParamOption struct {
	Name     string
	Keywords []string
	Desc     string
	Check    Validator
	Optional bool
}

func (opt ParamOption) Id() string {
	return opt.Name
}

func (opt ParamOption) Match() []string {
	return opt.Keywords
}

func (opt ParamOption) Describe() string {
	return opt.Desc
}

func (opt ParamOption) Validate(in string) (any, error) {
	if opt.Check == nil {
		panic("parsex.ParamOption{}.Check is nil")
	}

	return opt.Check(in)
}
