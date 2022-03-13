package io

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/jesperkha/Fizz/env"
)

// Standard io package for standard io operations

type i interface{}

var (
	scanner        = bufio.NewScanner(os.Stdin)
	ErrInvalidPath = errors.New("invalid path, line %d")
)

/*
	Gets user input from stdin.
	func input(prompt string) string
*/
func Input(prompt string) (input i, err error) {
	fmt.Print(prompt)
	scanner.Scan()
	return scanner.Text(), nil
}

/*
	Reads file and returns text.
	func readFile(filename string) string
*/
func ReadFile(filename string) (str i, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return str, ErrInvalidPath
	}

	return string(content), err
}

/*
	Writes content to file. Overwrites previous file content.
	func writeFile(filename string, content string)
*/
func WriteFile(filename string, content string) (val i, err error) {
	err = ioutil.WriteFile(filename, []byte(content), fs.ModeAppend)
	if err != nil {
		return val, ErrInvalidPath
	}

	return val, err
}

/*
	Appends content to file.
	func appendFile(filename string, content string)
*/
func AppendFile(filename string, content string) (val i, err error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.ModeAppend)
	if err != nil {
		return val, ErrInvalidPath
	}

	_, err = f.WriteString(content)
	if err != nil {
		return val, err
	}

	f.Close()
	return val, err
}

/*
	Returns list of files/directories in dir.
	func readDir(dir string) []string
*/
func ReadDir(dir string) (val i, err error) {
	dirs, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrInvalidPath
	}

	names := []interface{}{}
	for _, d := range dirs {
		names = append(names, d.Name())
	}

	return env.NewArray(names), err
}

/*
	Returns current working directory
	func curDir() string
*/
func CurDir() (str i, err error) {
	return os.Getwd()
}

/*
	Returns true if file exists.
	func exists(filename string) bool
*/
func Exists(filename string) (val i, err error) {
	_, e := os.Open(filename)
	return e == nil, err
}

/*
	Creates new directory.
	func newDir(name string)
*/
func NewDir(dirname string) (val i, err error) {
	os.Mkdir(dirname, os.ModeAppend)
	return nil, err
}

/*
	Creates new file. If the file already exists it will be overwritten.
	func newFile(name string)
*/
func NewFile(filename string) (val i, err error) {
	f, err := os.Create(filename)
	f.Close()
	return val, err
}
