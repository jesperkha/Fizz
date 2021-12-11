package interp

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/stmt"
	"github.com/jesperkha/Fizz/util"
)

var (
	ErrFileNotFound = errors.New("cannot find file with name: '%s'")
	ErrNonFizzFile  = errors.New("cannot run non-Fizz file")
)

// Interperates string of code. The string is tokenized in the lexer package
// where each token has a type, line, lexeme, and literal value. The parsed
// tokens are passed to the statement parser which looks for statements
// (identified by trailing semicolon or curly braces).

// The statements are parsed as statement tokens with a type, name, and
// expression values for the expression associated with the statement. For
// example: a variable declaration statement has an expression property for
// the initial value (if there is one) that is found after the equal sign.

// The statements are then executed and, when doing so, statement expressions
// are evaluated. Variable values are also assigned and manipulated in the
// variable environment found in the env package.

func Interperate(filename string, input string) (e env.Environment, err error) {
	// Parses input characters into lexical tokens for single and double symbols,
	// identifiers, and keywords.
	lexicalTokens, err := lexer.GetTokens(input)
	if err != nil {
		return e, err
	}

	// Lexical tokens are analysed and put into statement tokens. These statements
	// contain all the information they need for execution and error handling.
	statements, err := stmt.ParseStatements(lexicalTokens)
	if err != nil {
		return e, err
	}

	// File imports are handled after parsing the statements, not when executing them.
	// This means that all imports are run before anything else; they are "hoisted".
	// Even declaring a variable with the same name as the file before importing it will
	// raise an error as the file is imported before the variable is created.
	for _, s := range statements {
		if s.Type != stmt.Import {
			continue
		}
		
		e, err = RunFile(s.Name + ".fizz")
		if err != nil {
			return e, err
		}
		
		// Adds the global environment of the imported file to the env of the current one.
		// It is added as an object instance with the name file_import.
		split := strings.Split(s.Name, "/")
		if err = env.AddImportedFile(split[len(split)-1], e); err != nil {
			return e, err
		}
	}
	
	// Set origin point for function declarations. This makes sure that errors give
	// the correct filename when printed.
	stmt.CurrentOrigin = filename
	
	// Finally executes statement tokens. This is the only step that has any effect
	// on the actual input program as the others were just braking it up into usable
	// pieces. While the interpreter is still running, the values of variables will be
	// remembered as the environments are never reset at runtime.
	err = stmt.ExecuteStatements(statements)
	return env.NewEnvironment(), err
}

// Runs a fizz file. Is called from the run package upon running a script file, or
// in the Interperate() function, where imports are run as files and the environment
// is extracted and packaged into a namespace. Said namespace is put into the file
// it as imported to, which means if "main.fizz" imports "other.fizz", the main file
// also imports all of the files imported in "other.fizz".
func RunFile(filename string) (e env.Environment, err error) {
	if !strings.HasSuffix(filename, ".fizz") {
		return e, ErrNonFizzFile
	}

	if byt, err := os.ReadFile(filename); err == nil {
		e, err = Interperate(filename, string(byt))
		return e, util.WrapFilename(filename, err)
	}

	// Assumes path error
	return e, fmt.Errorf(ErrFileNotFound.Error(), filename)
}
