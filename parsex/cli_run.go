package parsex

import (
	"fmt"
	"os"
	"strings"
)

func (cli *CLI) simplifyArgs() *CLI {
	var newArgs = make([]string, 0, len(cli.args))
	for _, arg := range cli.args {
		newArgs = append(newArgs, strings.Split(arg, "=")...)
	}
	cli.args = newArgs

	return cli
}

func (cli *CLI) parse(inherit Input) (*CLI, Input, error) {
	output := inherit

	for i := 0; i < len(cli.args); i++ {
		if cli.args[i] == "--" {
			output.args = append(output.args, cli.args[i+1:]...)
			break
		}

		if strings.HasPrefix(cli.args[i], "-") {
			if err := cli.handleOption(cli.args[i], &output, &i); err != nil {
				return cli, output, err
			}
			continue
		}

		if branch, err := cli.tryBranch(cli.args[i], i); branch != nil {
			if err != nil {
				return cli, output, err
			}
			branch.parent = cli
			return branch.parse(output)
		}

		output.args = append(output.args, cli.args[i])
	}

	return cli, output, nil
}

func (cli *CLI) handleOption(arg string, output *Input, i *int) error {
	keyword := strings.TrimLeft(arg, "-")
	if index, ok := cli.keywords[keyword]; ok {
		switch opt := cli.opts[index].(type) {
		case ParamOption:
			if *i+1 >= len(cli.args) {
				return fmt.Errorf("%s: option %q requires a value", cli.Id(), arg)
			}
			value, err := opt.Validate(cli.args[*i+1])
			if err != nil {
				return fmt.Errorf("%s: option %q: %w", cli.Id(), arg, err)
			}
			output.values[opt.Id()] = value
			*i = *i + 1

		case FlagOption:
			output.values[opt.Id()] = true
		}

		return nil
	}

	for _, flag := range strings.Split(keyword, "") {
		index, ok := cli.keywords[flag]
		if !ok {
			return fmt.Errorf("%s: unknown option %q", cli.Id(), arg)
		}

		switch opt := cli.opts[index].(type) {
		case FlagOption:
			output.values[opt.Id()] = true

		default:
			return fmt.Errorf("%s: combined option %q can only be made out of flags", cli.Id(), arg)
		}
	}
	return nil
}

func (cli *CLI) tryBranch(arg string, idx int) (*CLI, error) {
	if index, ok := cli.keywords[arg]; ok {
		if branch, ok := cli.opts[index].(*CLI); ok {
			branch.args = cli.args[idx+1:]
			return branch, nil
		}
	}

	return nil, nil
}

func (cli *CLI) RunAndExit() {
	err := cli.Run()
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func (cli *CLI) Run() error {
	i := 0
	input := newInput()

	cli, input, err := cli.simplifyArgs().parse(newInput())
	if err != nil {
		return err
	}

	if input.Has("version") {
		cli.PrintVersion()
		return nil
	}

	if input.Has("help") {
		cli.PrintHelp()
		return nil
	}

	for _, arg := range cli.positional {
		if !strings.HasPrefix(arg, "?") {
			if i >= len(input.args) {
				return fmt.Errorf("%s: positional argument %q is missing", cli.Id(), arg)
			}
			i++
		}
	}

	for _, opt := range cli.opts {
		switch opt := opt.(type) {
		case ParamOption:
			if !opt.Optional {
				if !input.Has(opt.Id()) {
					return fmt.Errorf("%s: option %q is missing", cli.Id(), opt.Id())
				}
			}
		}
	}

	return cli.callback(input)
}
