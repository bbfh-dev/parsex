package parsex

import (
	"fmt"
)

type Builder struct {
	cli *CLI
}

func New(name string, desc string, callback Callback) *Builder {
	return &Builder{
		cli: &CLI{
			parent:     nil,
			Name:       name,
			Desc:       desc,
			callback:   callback,
			positional: []string{},
			opts: []Option{
				HelpFlag,
			},
			version:        "",
			args:           []string{},
			keywords:       map[string]int{},
			longestKeyword: 0,
		},
	}
}

func (builder *Builder) AddOptions(opts ...Option) *Builder {
	builder.cli.opts = append(builder.cli.opts, opts...)
	return builder
}

func (builder *Builder) AddArguments(args ...string) *Builder {
	builder.cli.positional = append(builder.cli.positional, args...)
	return builder
}

func (builder *Builder) SetVersion(version string) *Builder {
	builder.cli.version = version
	return builder.AddOptions(VersionFlag)
}

func (builder *Builder) Build() *CLI {
	for i, opt := range builder.cli.opts {
		for _, match := range opt.Match() {
			if _, ok := builder.cli.keywords[match]; ok {
				panic(
					fmt.Sprintf(
						"parsex.Builder.Build() found a keyword %q that was already previously used by another [parsex.Option]",
						match,
					),
				)
			}
			builder.cli.keywords[match] = i
		}
		if length := len(getMatches(opt)); builder.cli.longestKeyword < length {
			builder.cli.longestKeyword = length
		}
	}

	return builder.cli
}
