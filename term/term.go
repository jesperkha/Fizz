package term

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Terminal argument flag parser

var (
	//go:embed help.txt
	HELP            string
	ErrNotValidArg  = errors.New("flag name '%s' is not defined")
	ErrTooManyFiles = errors.New("cannot run more than one file")
)

// Arg passed is the value after an equals sign: --flag="arg". Defualt value is empty string
// Error returned is printed to stderr and further execution/parsing is stopped
type ArgFunc func()

type FlagParser struct {
	// Map of functions that run for specific flags
	Funcs map[string]ArgFunc
	// Bool flags for if a flag is present in args
	Flags map[string]bool
}

// Flag parser parses list of args in os.Args and converts them into map of bools
// Args are valid single-dash and double-dash flag names (without the dashes)
func NewFlagParser(sd []string, dd []string) *FlagParser {
	parser := FlagParser{
		map[string]ArgFunc{},
		map[string]bool{},
	}

	for _, arg := range sd {
		parser.Flags[fmt.Sprintf("-%s", arg)] = false
	}

	for _, arg := range dd {
		parser.Flags[fmt.Sprintf("--%s", arg)] = false
	}

	return &parser
}

func (p *FlagParser) IsValid(name string) bool {
	_, ok := p.Flags[name]
	return ok
}

func (p *FlagParser) IsFunc(name string) bool {
	_, ok := p.Funcs[name]
	return ok
}

// Maps function to flag name (including dash)
// Returns error if name was not defined in NewFlagParser()
func (p *FlagParser) Assign(name string, f ArgFunc) error {
	if !p.IsValid(name) {
		return fmt.Errorf(ErrNotValidArg.Error(), name)
	}

	p.Funcs[name] = f
	return nil
}

// Any undefined flag names found will return an error containing the first invalid flag name
func (p *FlagParser) Parse() (filename string, err error) {
	args := os.Args[1:]

	for _, flag := range args {
		if strings.HasPrefix(flag, "--") {
			if !p.IsValid(flag) {
				return filename, fmt.Errorf(ErrNotValidArg.Error(), flag)
			}

			if p.IsFunc(flag) {
				p.Funcs[flag]()
			}

			p.Flags[flag] = true
			continue
		}

		if strings.HasPrefix(flag, "-") {
			for _, opt := range strings.Split(flag, "-")[1] {
				opt_s := "-" + string(opt)
				if !p.IsValid(opt_s) {
					return filename, fmt.Errorf(ErrNotValidArg.Error(), opt_s)
				}

				if p.IsFunc(opt_s) {
					p.Funcs[opt_s]()
				}

				p.Flags[opt_s] = true
			}

			continue
		}

		if filename != "" {
			return filename, ErrTooManyFiles
		}

		filename = flag
	}

	return filename, err
}
