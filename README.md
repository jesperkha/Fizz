# Fizz

## **About**

Fizz is an interpreted programming language built with Go. The main use cases for Fizz is file management, simple http requests, and general terminal scripting. The built in libraries (not added yet) support all of these features, including documentation with examples.

<br>

## **Roadmap**

This roadmap highlights progress for Fizz's development:

- ✔️ Expression parsing
- ✔️ Conditional statements
- ✔️ Loops
- ❌ Functions
- ❌ Classes
- ❌ Arrays
- ❌ File import
- ❌ Go -> Fizz library support
- ❌ Language documentation

<br>

## **Download / Setup**

> **Disclaimer**: You need to have Go installed on your device to build the executable or run the interpreter.

<br>

First clone the repository to your device. To run the interpreter either build an executable with `build.bat` or run `go run . [filename]`

<br>

### **Terminal mode**

Running the interpreter without giving a filename with run the terminal mode where you can run any valid Fizz code live. Errors are printed but the program is not terminated.

<br>

### **Run file**

Running the interpreter and giving a filename simply runs the code in the file and halts if an error occurs. Fizz files must end in the `.fizz` suffix.

```c
fizz myFile.fizz
fizz myFile // valid if the file ends with .fizz
```
