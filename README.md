# Parsex

A GNU-/POSIX-compiant CLI argument parsing and validation library.

### Features

- Supports flags, options and subcommands.
- Recognizes all argument formats: `-a`, `-abc`, `-flag`, `-opt=value`, `-opt value`, `--flag`, `--flag=value`, `--flag value`.
- Supports `--` to separate arguments.

### Table of contents

<!-- vim-markdown-toc GFM -->

- [Example Usage](#example-usage)

<!-- vim-markdown-toc -->

# Example Usage

```go
package main

import "github.com/bbfh-dev/parsex/parsex"

// Use Builder pattern to construct your CLI
var CLI = parsex.
    New(
        // The name of the executable / branch
        "example",
        // Description to show up in --help
        "An example description for the command",
        // Callback function (using anonymous for demonstration, you would usually define it somewhere else)
        func(in parsex.Input) error {
            var length int = in.Int("length")  // Example of how you can access integers
            var debug bool = in.Has("debug")  // Example of checking flags / if options are provided

            return nil
        },
    ).
    // Adds --version support to the CLI
    SetVersion("v1.0.0-beta").
    // These are flags, options and subcommands
    AddOptions(
        parsex.FlagOption{
            Name: "debug",
            // --debug, -debug, -D, --D will be recognized
            Keywords: []string{"debug", "D"},
            Desc:     "Debug this",
        },
        // You can create branches (subcommands) from other CLIs.
        parsex.New("subcommand", "This is my subcommand!", SubcommandFunc).Build(),
        parsex.ParamOption{
            Name:     "input",
            Keywords: []string{"input", "i"},
            Desc:     "Input file",
            Check:    parsex.ValidPath,
            Optional: true,
        },
    ).
    // Inform user about positional arguments
    // Prefix with "?" if it's optional. Otherwise parsex will force user to provide it.
    AddArguments(
        "?optional",
        "required",
        "more...",
    ).
    // We're done here, get the CLI.
    Build()

func main() {
    // Use os.Args as input, run the program and handle errors automatically
    CLI.FromOSArgs().RunAndExit()
}
```
