package io

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

// Standard io package for standard io operations

type i interface{}

var (
	Includes = map[string]interface{}{}
	scanner  = bufio.NewScanner(os.Stdin)
)

func init() {
	Includes = map[string]interface{}{
		// Get user input from stdin
		"input": input,
		// Read file and return string
		"readFile": readFile,
		// Write content to file
		"writeFile": writeFile,
		// Append content to file
		"appendFile": appendFile,
	}
}

func input(prompt string) (input i, err error) {
	fmt.Print(prompt)
	scanner.Scan()
	return scanner.Text(), nil
}

func readFile(filename string) (str i, err error) {
	content, err := ioutil.ReadFile(filename)
	return string(content), err
}

func writeFile(filename string, content string) (val i, err error) {
	err = ioutil.WriteFile(filename, []byte(content), fs.ModeAppend)
	return val, err
}

func appendFile(filename string, content string) (val i, err error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.ModeAppend)
	if err != nil {
		return val, err
	}

	_, err = f.WriteString(content)
	if err != nil {
		return val, err
	}

	f.Close()
	return val, err
}

// read dir, curdir, 