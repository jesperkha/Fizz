<br />
<div align="center">
  <img src=".github/logo.svg" alt="Logo" width="120">

  <h3 align="center">Fizz, the programming language</h3>

  <p align="center">
    Interpreted dynamic programming language built with Go.
    <br />
    <a href="https://github.com/jesperkha/Fizz/blob/main/docs/lang.md"><strong>Documentation »</strong></a>
    <br />
    <br />
    <a href="https://github.com/jesperkha/Fizz/tree/main/examples">Examples</a>
    ·
    <a href="https://github.com/jesperkha/Fizz/issues">Report Bug</a>
    ·
    <a href="#installation">Download</a>
  </p>
</div>

<br>

<details>
  <summary>Table of Contents</summary>
  <ul>
  <li><a href="#about">About</a></li>
  <li><a href="#installation">Installation</a></li>
  <li><a href="#documentation">Documentation</a></li>
  <li><a href="#running-a-program">Running a program</a></li>
  <li><a href="#code-examples">Code examples</a></li>
  </ul>
</details>

<br>

## About

Fizz is a dynamic and interpreted programming language built with Go. It is strongly typed and comes with very readable and accurate error messages. Fizz has most of the standard functionality that you would expect from modern programming languages. The library system also allows the user to implement their own features as Go functions and port them directly into Fizz. If you like this project, consider giving it a star 😉

### Features

- Variables, conditionals, and loops
- Functions, arrays, and objects
- File imports and libraries
- Clean syntax and simple grammar

<br>

## Installation

Prebuilt binary of the [latest release (v1.1.0)](https://github.com/jesperkha/Fizz/releases/tag/v1.1.0).

You can also build from source. However, building from source from a non-release branch does not guarantee that everything works as expected as some things may be undergoing changes.

1. Clone repo
2. Run the `build.sh` file

### Syntax highlighting

Finally, there is also optional, but recommended, [syntax highlighting](https://github.com/jesperkha/fizz-extensions) extensions for both Visual Studio Code and micro.

<br>

## Documentation

You can read the [full language documentation](./docs/lang.md) to learn about all of Fizz's syntax. It is also recommended to quickly skim over [the language grammar](./docs/grammar.md) to make sure you undestand the basics of how Fizz is structured (don't worry, it's _very_ similar to most modern programming languages).

Make sure to check out [the command-line basics](./docs/cmd.md) too so you know how to run your code and also which configurations you can apply.

<br>

## Running a program

[Full documentation on command-line basics](./docs/cmd.md)

### Terminal mode

Running the interpreter without giving a filename will run the terminal mode where you can run any valid Fizz code live. Errors are printed but the program is not terminated. Newlines are supported for blocks and the code will not be executed until the block is closed.

### Run file

Running the interpreter and giving a filename simply runs the code in the file and halts if an error occurs. Fizz files must end with the `.fizz` suffix. Both of the following are valid:

```console
$ ./fizz myFile.fizz
$ ./fizz myFile
```

<br>

## Code examples

Some simple code examples written in Fizz. There are more and bigger examples in the `examples` directory. All of the features used here and many more are thoroughly documented in the [documentation page](./docs/lang.md).

<br>

Write to a file:

```go
include "io";
include "str";

define Person {
  name, age
}

func writePerson(person) {
  if person == nil : person.name == "" {
    error "Please enter valid person";
  }

  io.appendFile("names.txt", str.format(person));
}

john := Person("John", 59);
writeName(john);
```

<br>

Find max and min numbers in array:

```go
include "str";

arr := [5, 3, 7.5, 8, 2];
max := 0;
min := 999;

range n in arr {
  if n > max {
    max = n;
  }
  if n < min {
    min = n;
  }
}

print "Max: " + str.toString(max);
print "Min: " + str.toString(min);
```
