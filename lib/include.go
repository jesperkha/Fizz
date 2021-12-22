package lib

import "github.com/jesperkha/Fizz/lib/std"

func init() {
	Add("std", std.Get())
}