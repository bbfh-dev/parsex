# Parsex

A GNU-/POSIX-compiant CLI argument parsing and validation library.

<!-- vim-markdown-toc GFM -->

* [Argument types](#argument-types)
    * [Short option flags](#short-option-flags)
    * [Short options with arguments](#short-options-with-arguments)
    * [Long option flags](#long-option-flags)
    * [Long options with arguments](#long-options-with-arguments)
* [Positional arguments](#positional-arguments)
* [Development](#development)

<!-- vim-markdown-toc -->

# Argument types

## Short option flags

**Example**: `-a -b -c` / `-abc`

**Usage**: Simple boolean flags

## Short options with arguments

**Example**: `-A/tmp/file.png` / `-A /tmp/file.png`

**Usage**: Simple parameters with a value. If no space is provided the argument must be capitalized.

## Long option flags

**Example**: `--version`

**Usage**: Same as [Short option flags](#short-option-flags), but taken from GNU

## Long options with arguments

**Example**: `--output=/tmp/file.png` / `--output /tmp/file.png`

**Usage**: Same as [Short options with arguments](#short-options-with-arguments), but taken from GNU. Can be separated with either space or an equal sign (=).


# Positional arguments

Any arguments subsequent to options/flags will be parsed as positional arguments.

Using `--` causes everything to the right to be threated as positional arguments.

# Development

- Use `make test` to run tests on the project.
- Use `make coverage` to create tests/coverage.html with the test coverage of the library.
