package parsex

import (
	"fmt"
	"os"
	"strings"
)

type Callback func(Input) error

type CLI struct {
	parent *CLI

	Name string
	Desc string

	callback       Callback
	positional     []string
	opts           []Option
	keywords       map[string]int
	longestKeyword int

	version string
	args    []string
}

func (cli *CLI) Id() string {
	return cli.Name
}

func (cli *CLI) Match() []string {
	return []string{cli.Name}
}

func (cli *CLI) Describe() string {
	return cli.Desc
}

func (cli *CLI) FromSlice(in []string) *CLI {
	cli.args = in
	return cli
}

func (cli *CLI) FromOSArgs() *CLI {
	if len(os.Args) < 2 {
		return cli
	}
	cli.args = os.Args[1:]
	return cli
}

func (cli *CLI) printFullName() string {
	if cli.parent == nil {
		return cli.Id()
	}
	return cli.parent.printFullName() + " " + cli.Id()
}

func (cli *CLI) PrintVersion() {
	if cli.version == "" {
		fmt.Println(cli.printFullName())
		return
	}
	fmt.Println(cli.printFullName(), cli.version)
}

func (cli *CLI) printHelpUsage() {
	fmt.Printf("\nUsage:\n\t%s <options> ", cli.Id())
	for _, arg := range cli.positional {
		if strings.HasPrefix(arg, "?") {
			fmt.Printf("(%s) ", arg[1:])
		} else {
			fmt.Printf("[%s] ", arg)
		}
	}
	fmt.Print("\n")
}

func (cli *CLI) printHelpCommands() {
	fmt.Print("\nCommands:")
	for _, opt := range cli.opts {
		switch opt := opt.(type) {
		case *CLI:
			fmt.Printf("\n\t%s\t%s", opt.Id(), opt.Describe())
		}
	}
	fmt.Print("\n")
}

func (cli *CLI) printHelpOptions() {
	fmt.Print("\nOptions:")
	for _, opt := range cli.opts {
		switch opt := opt.(type) {
		case FlagOption:
			matches := getMatches(opt)
			fmt.Printf("\n\t%s  %s  %s", matches, strings.Repeat(" ", max(0, cli.longestKeyword-len(matches))), opt.Describe())
		}
	}
	fmt.Print("\n")
}

func (cli *CLI) PrintHelp() {
	cli.PrintVersion()
	cli.printHelpUsage()
	cli.printHelpCommands()
	cli.printHelpOptions()
}
