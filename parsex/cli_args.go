package parsex

import (
	"fmt"
	"strings"
	"unicode"
)

type Arg struct {
	// Name is used for you to identify this argument later inside of parsex.Input
	Name string

	// Comma separated strings to match inside of input arguments.
	//
	// Example: "--output,-o" or "--AUTO,-o" where Arg.Name = "output".
	//
	// any occurance of `--AUTO` will be replaced with the lowecase Arg.Name
	Match string

	// Description of the argument to be printed in the help documentation
	Desc string

	// (Optional) Require an argument after the option of a specific valid value
	Check Validator

	// (Optional) Make this argument a subcommand, all arguments aftwards will be handled by this Arg.Branch
	Branch *CLI
}

func applyAutoFlags(args []Arg) []Arg {
	for i := range args {
		args[i].Match = strings.ReplaceAll(args[i].Match, "AUTO", strings.ToLower(args[i].Name))
	}
	return args
}

// Creates a map of singular options to their coresponding arg.Match
func argToIndexMap(args []Arg) map[string]int {
	out := map[string]int{}

	for i, arg := range args {
		for _, key := range strings.Split(arg.Match, ",") {
			out[key] = i
		}
	}

	return out
}

func isCombinedArg(arg string) bool {
	return strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") &&
		!unicode.IsUpper(rune(arg[1]))
}

func simplifyArgs(args []string) []string {
	var out []string

	for _, arg := range args {
		if len(arg) == 0 {
			continue
		}

		if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") {
			out = append(out, strings.Split(arg, "=")...)
			continue
		}

		if strings.HasPrefix(arg, "-") && unicode.IsUpper(rune(arg[1])) && len(arg) > 2 {
			out = append(out, arg[:2], arg[3:])
			continue
		}

		if isCombinedArg(arg) {
			for _, char := range arg[1:] {
				out = append(out, fmt.Sprintf("-%s", string(char)))
			}
			continue
		}

		out = append(out, arg)
	}

	return out
}
