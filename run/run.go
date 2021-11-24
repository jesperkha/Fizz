package run

import (
	"errors"
	"fmt"
	"strings"
)

// Todo: Move this to main

var (
	ErrFileNotFound = errors.New("cannot find file with name: '%s'")
	ErrNonFizzFile  = errors.New("cannot run non-Fizz file")
)

var cmdOptions = map[string]func(){
	"help": func() {
		fmt.Println("use: fizz [filename.fizz | --option]")
	},
}

func RunInterpeter(args []string) {
	if len(args) == 2 {
		arg := args[1]

		// Run option commands
		if strings.HasPrefix(arg, "--") {
			if opt, ok := cmdOptions[strings.TrimLeft(arg, "-")]; ok {
				opt()
				return
			}

			fmt.Printf("unknown option: '%s'\n", arg)
			return
		}

		// Run fizz file
		filename := args[1]
		if !strings.HasSuffix(filename, ".fizz") {
			filename += ".fizz"
		}

		if err := RunFile(filename); err != nil {
			fmt.Println(err)
		}

		return
	}

	// Run terminal mode
	RunTerminal()
}
