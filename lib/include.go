package lib

import "github.com/jesperkha/Fizz/lib/std"

// Todo: make dynamic import of folders
func init() {
	Add("std", FuncMap{
		"input": std.GetStdinInput,
	})
}