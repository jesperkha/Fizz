package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/interp"
	"github.com/jesperkha/Fizz/lib"
	"github.com/jesperkha/Fizz/stmt"
	"github.com/jesperkha/Fizz/term"
	"github.com/jesperkha/Fizz/util"
)

var (
	ErrOneArgOnly = errors.New("expected a single argument, got %d")
	validArgs     = []string{"--help", "--version", "-f", "-e"}
)

func RunInterpreter() {
	parser, err := term.Parse(validArgs)
	if err != nil {
		util.ErrorAndExit(err)
	}
	args := parser.Args()
	if len(args) > 1 {
		util.ErrorAndExit(fmt.Errorf(ErrOneArgOnly.Error(), len(args)))
	}

	// Early exit options
	if parser.HasOption("help") {
		fmt.Println(term.HELP)
		return
	} else if parser.HasOption("version") {
		fmt.Printf("Fizz %s\n", VERSION)
		return
	}

	// Subcommands
	switch parser.SubCommand() {
	case "docs":
		if err := lib.PrintDocs(args[0]); err != nil {
			util.ErrorAndExit(err)
		}
		return
	case "help":
		if msg, ok := term.CommandDescriptions[args[0]]; ok {
			fmt.Println(msg)
		} else {
			util.PrintError(fmt.Errorf(term.ErrUnknownCommand.Error(), args[0]))
		}
		return
	}

	// Run terminal mode if no other args are given
	if len(args) == 0 {
		RunTerminal()
		return
	}

	// Goto directory of file specified
	split := strings.Split(args[0], "/")
	path := strings.Join(split[:len(split)-1], "/")
	name := split[len(split)-1]
	os.Chdir(path)

	// Run file
	e, err := interp.RunFile(name)

	// Print global environment if flag is set first
	if parser.HasFlag("e") {
		fmt.Println(util.FormatPrintValue(e))
	}

	// Handle error
	if err != nil && err != stmt.ErrProgramExit {
		util.PrintError(err)
		if c := env.GetCallstack(); parser.HasFlag("f") && len(c) > 0 {
			util.PrintError(fmt.Errorf(c))
		}

		os.Exit(1)
	}
}

// Leaves the interpreter running as the user inputs code to the terminal.
// Prints out errors but does not terminate until ^C or 'exit'.
func RunTerminal() {
	fmt.Println("type 'exit' to terminate session")
	scanner := bufio.NewScanner(os.Stdin)
	numBlocks, line := 0, 1
	totalString, space := "", " "
	env.ThrowEnvironment = false

	for {
		fmt.Printf("%d%s : %s", line, space, strings.Repeat("    ", numBlocks))
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		// Continue with indent after braces
		numBlocks += strings.Count(input, "{") - strings.Count(input, "}")
		totalString += input + "\n" // Better error handling
		if numBlocks <= 0 {
			if _, err := interp.Interperate("", totalString); err != nil {
				util.PrintError(err)
				line--
			}

			totalString = ""
			numBlocks = 0
		}

		line++
		if line == 10 {
			space = ""
		}
	}

	fmt.Println("session ended")
}
