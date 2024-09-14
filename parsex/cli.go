package parsex

import (
	"fmt"
	"os"
	"strings"

	"github.com/bbfh-dev/parsex/parsex/internal"
)

// Function to be run with the provided arguments and flags
type Program func(in Input, args ...string) error

// The primary handler for a single sub-command
type CLI struct {
	Name       string
	cliProgram Program
	cliArgs    []Arg
	argToIndex map[string]int
	arguments  []string
	input      Input
	inArgs     []string
}

// Create a new program. Use for both main entry and branches (sub-commands)
func New(name string, program Program, args []Arg) *CLI {
	prepend := []Arg{{Name: "help", Match: "--AUTO", Desc: "print help message and exit"}}
	args = applyAutoFlags(append(prepend, args...))
	return &CLI{
		Name:       name,
		cliProgram: program,
		cliArgs:    args,
		argToIndex: argToIndexMap(args),
		arguments:  []string{},
		input:      Input{},
		inArgs:     []string{},
	}
}

func (cli *CLI) reset() *CLI {
	cli.input = Input{}
	cli.inArgs = []string{}

	return cli
}

// Set input arguments from a string
//
// Example: `-o /tmp/ --verbose`
func (cli *CLI) FromString(str string) *CLI {
	cli.arguments = simplifyArgs(strings.Split(str, " "))
	return cli.reset()
}

// Set input arguments from os.Args (provided from calling the binary)
func (cli *CLI) FromArgs() *CLI {
	if len(os.Args) > 1 {
		cli.arguments = simplifyArgs(os.Args[1:])
	}
	return cli.reset()
}

// Set input arguments directly
func (cli *CLI) FromSlice(args []string) *CLI {
	cli.arguments = simplifyArgs(args)
	return cli.reset()
}

// Parse input arguments and execute the correct program
func (cli *CLI) Run() error {
	var expectsArgumentsOnly bool
	var i int
	for i < len(cli.arguments) {
		if expectsArgumentsOnly && strings.HasPrefix(cli.arguments[i], "-") {
			return fmt.Errorf(
				"Unexpected option %q after arguments. Make sure all arguments are provided in the end or after --",
				cli.arguments[i],
			)
		}

		if cli.arguments[i] == "--" {
			if len(cli.arguments) > i+1 {
				cli.inArgs = append(cli.inArgs, cli.arguments[i+1:]...)
			}
			break
		}

		arg, err := cli.findArg(cli.arguments[i], i)
		if err != nil {
			if !strings.HasPrefix(cli.arguments[i], "-") {
				cli.inArgs = append(cli.inArgs, cli.arguments[i])
				expectsArgumentsOnly = true
				i += 1
				continue
			}
			return err
		}

		if arg.Branch != nil {
			arg.Branch.input = cli.input
			arg.Branch.inArgs = cli.inArgs
			if len(cli.arguments) > i+1 {
				arg.Branch.arguments = cli.arguments[i+1:]
			}
			return internal.PrefixErr(arg.Name, arg.Branch.Run())
		}

		diff, err := cli.processArg(arg, cli.arguments[i:])
		if err != nil {
			return err
		}

		i += diff
	}

	if cli.input.Has("help") {
		fmt.Println(cli.Help(cli.Name))
		return nil
	}

	return cli.cliProgram(cli.input, cli.inArgs...)
}

func (cli *CLI) findArg(current string, i int) (Arg, error) {
	for keyFrom := range cli.argToIndex {
		if keyFrom == current {
			return cli.cliArgs[cli.argToIndex[current]], nil
		}
	}

	return Arg{}, fmt.Errorf("Unknown argument (#%d): %s", i+1, current)
}

func (cli *CLI) processArg(arg Arg, source []string) (int, error) {
	if arg.Check == nil {
		cli.input[arg.Name] = ""
		return 1, nil
	}

	if len(source) < 2 {
		return 1, fmt.Errorf("Argument %q requires a value!", arg.Name)
	}

	cli.input[arg.Name] = source[1]

	return 2, nil
}
