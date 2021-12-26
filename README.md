# **Fizz**

Fizz is an interpreted programming language built with Go. The main use cases for Fizz will be file management, simple http requests, and general terminal scripting. The built in libraries (not added yet) will support all of these features, including documentation with examples.

<br>

## Table of contents

- [Roadmap](#roadmap)
- [Language documentation](#language-documentation)
- [Installation](#installation)
  - [Pre-built binaries](#pre-built-binaries)
  - [Building from source](#building-from-source)
  - [Extensions](#extensions)
- [Running a program](#running-a-program)

<br>

---

<br>

## <a id="roadmap"></a> Roadmap

This roadmap highlights progress for the development of Fizz:

- ✔️ Expression parsing
- ✔️ Conditional statements
- ✔️ Loops
- ✔️ Functions
- ✔️ Objects
- ✔️ File import
- ✔️ Go -> Fizz library support
- ✔️ Arrays
- ❌ Hashtables
- ❌ Full language documentation

<br>

## <a id="language-documentation"></a> Language documentation

You can find all the basic info you need about Fizz [here](./docs/lang.md). For a deeper understanding of how Fizz works its suggested to simply look at the source code.

<br>

## <a id="installation"></a> Installation

<a id="pre-built-binaries"></a> Pre-built binaries:

- [Latest Release v0.5.0](https://github.com/jesperkha/Fizz/releases/tag/v0.5.0)

<a id="building-from-source"></a> Build from source:

1. Clone repo
2. Run the `build.sh` file

<a id="extensions"></a> Extensions:

- [Syntax highlighting](https://github.com/jesperkha/fizz-extensions)

<br>

## <a id="running-a-program"></a> Running a program

**Terminal mode**

Running the interpreter without giving a filename will run the terminal mode where you can run any valid Fizz code live. Errors are printed but the program is not terminated. Newlines are supported for blocks and the code will not be executed until the block is closed.

**Run file**

Running the interpreter and giving a filename simply runs the code in the file and halts if an error occurs. Fizz files must end in the `.fizz` suffix.

Both of these are valid:

```console
$ ./fizz myFile.fizz
$ ./fizz myFile
```
