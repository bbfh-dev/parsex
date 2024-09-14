package parsex

import "strings"

const HELP_ARG_WIDTH = 16

func (cli *CLI) Help(cmd string) string {
	var builder strings.Builder

	builder.WriteString("Usage: " + cmd + " [OPTIONS] <ARGUMENTS...>\n")

	var commandsBuilder strings.Builder
	var optionsBuilder strings.Builder

	for _, arg := range cli.cliArgs {
		if arg.Branch == nil {
			optionsBuilder.WriteString("    " + arg.Match)
			optionsBuilder.WriteString(strings.Repeat(" ", HELP_ARG_WIDTH-len(arg.Match)))
			optionsBuilder.WriteString(arg.Desc + "\n")
		} else {
			commandsBuilder.WriteString("    " + arg.Match)
			commandsBuilder.WriteString(strings.Repeat(" ", HELP_ARG_WIDTH-len(arg.Match)))
			commandsBuilder.WriteString(arg.Desc + "\n")
		}
	}

	if len(commandsBuilder.String()) > 0 {
		builder.WriteString("\nCommands:\n")
		builder.WriteString(commandsBuilder.String())
	}

	if len(optionsBuilder.String()) > 0 {
		builder.WriteString("\nOptions:\n")
		builder.WriteString(optionsBuilder.String())
	}

	return builder.String()
}
