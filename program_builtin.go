package parsex

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

var indent = strings.Repeat(" ", 4)

var helpOption = option{
	Name: "help",
	Alt:  "",
	Desc: "Print this help message",
	Type: reflect.TypeOf(true),
	Ref:  nil,
}

var versionOption = option{
	Name: "version",
	Alt:  "",
	Desc: "Print the program version",
	Type: reflect.TypeOf(true),
	Ref:  nil,
}

func (program *Program) printVersion(writer io.Writer) {
	if program.version == "" {
		fmt.Fprintln(writer, program.Name)
		return
	}
	fmt.Fprintf(writer, "%s v%s\n", program.Name, program.version)
}

func (program *Program) printArgs(writer io.Writer) {
	for _, arg := range program.posArgs {
		writer.Write([]byte("<" + arg + "> "))
	}
}

func (program *Program) printHelp(writer io.Writer) {
	program.printVersion(writer)
	fmt.Fprintf(writer, "\n%s\n\nUsage:\n%s%s [options] ", program.Desc, indent, program.Name)
	program.printArgs(writer)
	fmt.Fprintf(writer, "\n")

	if len(program.branches) != 0 {
		fmt.Fprint(writer, "\nCommands:\n")
		for _, branchKey := range program.branchKeys {
			fmt.Fprintf(writer, "%s%s [options] ", indent, branchKey)
			program.branches[branchKey].printArgs(writer)
			fmt.Fprint(writer, "\n")
		}
	}

	if len(program.optionKeys) != 0 {
		fmt.Fprint(writer, "\nOptions:\n")
		for _, optionKey := range program.optionKeys {
			fmt.Fprintf(
				writer,
				"%s%s\n%s# %s\n",
				indent,
				program.options[optionKey],
				indent+indent,
				program.options[optionKey].Desc,
			)
		}
	}
}
