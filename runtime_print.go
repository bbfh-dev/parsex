package parsex

import (
	"fmt"
	"io"
	"strings"

	"github.com/bbfh-dev/parsex/internal"
)

var indent = strings.Repeat(" ", 4)

func (runtime *runtimeType) PrintVersion(writer io.Writer) {
	if runtime.version == "" {
		fmt.Fprintln(writer, runtime.name)
		return
	}
	fmt.Fprintf(writer, "%s v%s\n", runtime.name, runtime.version)
}

// Performs preprocessing and prints the help text block
func (runtime *runtimeType) SafePrintHelp(writer io.Writer) error {
	if err := runtime.preprocess(); err != nil {
		return err
	}
	runtime.printHelp(writer)
	return nil
}

func (runtime *runtimeType) printHelp(writer io.Writer) {
	runtime.PrintVersion(writer)
	fmt.Fprintf(writer, "\n%s\n\nUsage:\n%s%s [options] ", runtime.desc, indent, runtime.name)
	runtime.printArgs(writer)
	fmt.Fprintf(writer, "\n")

	if !runtime.branches.IsEmpty() {
		fmt.Fprint(writer, "\nCommands:\n")

		runtime.branches.ForEach(func(name string, branch *runtimeType) {
			fmt.Fprintf(writer, "%s%s [options] ", indent, name)
			branch.printArgs(writer)
			fmt.Fprint(writer, "\n")
		})
	}

	if !runtime.genOptions.IsEmpty() {
		fmt.Fprint(writer, "\nOptions:\n")

		runtime.genOptions.ForEach(func(_ string, option internal.Option) {
			fmt.Fprintf(
				writer,
				"%s%s\n%s# %s\n",
				indent,
				option.String(),
				indent+indent,
				option.Desc,
			)
		})
	}
}

func (runtime *runtimeType) printArgs(writer io.Writer) {
	for _, arg := range runtime.posArgs {
		writer.Write([]byte("<" + arg + "> "))
	}
}
