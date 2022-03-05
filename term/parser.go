package term

import (
	_ "embed"
	"os"
	"strings"
)

// Terminal argument parser

var (
	//go:embed help.txt
	HELP string
)

type ArgList interface {
	// Returns true if flag is present (argument starting with '-')
	HasFlag(flag string) bool

	// Returns true if option is present (argument starting with '--')
	HasOption(option string) bool

	// Returns the name of the subcommand used. The subcommand is the first
	// string found after the program name, unless it is the only argument,
	// in which case it will be handled as an argument, not a subommand.
	SubCommand() string

	// Returns the arguments, not flags, options, or the subcommand.
	Args() []string
}

// Parses arguments into ArgList. Raises error if an unknown flag, option, or
// subcommand is found. Todo: make flag validator (pass in valid arg list)
func Parse() (list ArgList, err error) {
	args := os.Args[1:]
	handler := ArgHandler{}

	for idx, arg := range args {
		if strings.HasPrefix(arg, "--") {
			// Check if option
			handler.options = append(handler.options, strings.TrimLeft(arg, "-"))
		} else if strings.HasPrefix(arg, "-") {
			// Check if flag (after to avoid false positive)
			handler.flags = append(handler.flags, strings.TrimLeft(arg, "-"))
		} else if idx == 0 && len(args) != 1 {
			// Check for valid subcommand
			handler.subcmd = arg
		} else {
			handler.args = append(handler.args, arg)
		}
	}

	return &handler, err
}
