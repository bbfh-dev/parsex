# Parsex

Parsex `/pɑːrs ˈɛks/` — a GNU-/POSIX-compiant CLI argument parsing and validation library.

## Features

- Supports `--flags`, `--options <value>` and `subcommands`.
- Recognizes all argument formats: `-a`, `-abc`, `-flag`, `-opt=value`, `-opt value`, `--flag`, `--flag=value`, `--flag value`.
- Supports `--` to separate arguments.

## Table of contents

<!-- vim-markdown-toc GFM -->

* [Example usage](#example-usage)

<!-- vim-markdown-toc -->

## Example usage

```go
package main

import (
    "os"

    "github.com/bbfh-dev/parsex"
)

var Options struct {
    Verbose    bool   `alt:"v" desc:"Print verbose debug information"`
    Debug      bool   `alt:"d" desc:"Run the program in DEBUG mode"`
    Input      string `desc:"Some input file"`
    SomeNumber int    `alt:"N" desc:"A valid integer"`
}

var program = parsex.Program{
    // You can set it to `nil` if you don't have any options
    Data: &Options,
    Name: "example",
    Desc: "This is an example program",
    Exec: func(args []string) error {
        // I use an annonymous function just for demo.
        // This is what's gonna be called with the positional arguments
        // provided in [args] and all options saved to [Options]
        return nil
    },
}.Runtime().SetVersion("1.0.0-dev").SetPosArgs("arg1", "arg2?", "argN...")

func main() {
    err := program.Run(os.Args)
    if err != nil {
        // All errors are typed, allowing you to know exactly what went wrong.
        // Regarless of how you handle them
        // err.Error() should be enough information for the user.
        switch err := err.(type) {
        case parsex.ErrExecution:
            // You can access the various properties provided into the error
            if err.ErrKind == parsex.ErrKindExecIsNil && err.Name == "example" {
                // In case you want to handle a specific error
            }
        default:
            os.Stderr.WriteString(err.Error() + "\n")
        }
    }
}
```
