# Fizz

## **About**

Fizz is an interpreted programming language built with Go. The main use cases for Fizz is file management, simple http requests, and general terminal scripting. The built in libraries (not added yet) support all of these features, including documentation with examples.

<br>

## **Roadmap**

This roadmap highlights progress for Fizz's development:

- ✔️ Expression parsing
- ✔️ Conditional statements
- ✔️ Loops
- ✔️ Functions
- ❌ Classes
- ❌ Arrays
- ❌ File import
- ❌ Go -> Fizz library support
- ❌ Language documentation

<br>

## **Language**

You can find all the basic info you need about Fizz [here](lang.md). For a deeper understanding of how Fizz works its suggested to simply look at the source code.

<br>

## **Setup and use**

#### Pre-built binaries:

- [Test Release](https://github.com/jesperkha/Fizz/releases/tag/test-release)

#### Or you can clone the repo and:

- build an executable with `build.bat`
- or run with `go run . [filename]`

<br>

### **Terminal mode**

Running the interpreter without giving a filename will run the terminal mode where you can run any valid Fizz code live. Errors are printed but the program is not terminated.

### **Run file**

Running the interpreter and giving a filename simply runs the code in the file and halts if an error occurs. Fizz files must end in the `.fizz` suffix.

Both of these are valid:

`fizz myFile.fizz`

`fizz myFile`
