# **Fizz libraries**

- [Overview](#overview)
- [Make your own](#create)
  - [Setup](#setup)
  - [Content](#content)
  - [Build](#build)
  - [Docs](#docs)

<br>

## <a id="overview"></a> **Overview**

Fizz allows for Go libraries to be included in your Fizz code. Use the `include` statement followed by the library name to use it:

```go
include "io";

name := io.input("Enter name: ");
```

<br>

## <a id="create"></a> **Making your own library**

You can easily create your own library for Fizz. The only requirements are:

- A unique library name as duplicates are not allowed
- That you follow the steps below to make sure your library is valid

<br>

### <a id="setup"></a> Setup

Create a new folder in the `lib` directory and name it after your library. The main file in your library **must** have the same name as the folder its in. Here is an example structure:

```
lib/
    mylib/
        mylib.go
```

<br>

### <a id="content"></a> Content

To get started, add the following to your main file:

```go
package mylib

var Includes = map[string]interface{}{}

func init() {

}
```

This is technically all you need to make a valid library. The name and type of `Includes` **must** be as shown. The `init()` function is where you add the functions you want to include in your library.

```go
package mylib

var Includes = map[string]interface{}{}

func init() {
    // Function will be exported with name "hello". Names do not need to match.
    Includes["hello"] = sayHello
}

// All functions in this package must return a value (interface) and error.
// If the error returned is not nil, it will be raised as a fizz error and
// terminate the program.
func sayHello(name string) (val interface{}, err error) {
    return "Hello, " + name, err
}
```

The values of the arguments given are checked before trying to call the function, so if the types do not match an error is raised. The return value will always be as shown.

<br>

## <a id="build"></a> Building with the new library

To actually add your new library to Fizz you need to recompile with the `build.sh` file. This will run a python script which will add an import to your library package.

When running Fizz again after compiling you can use your new library:

```go
include "mylib";

print mylib.hello("John");
```

```console
$ ls
main.fizz

$ fizz main
Hello, John
```

<br>

## <a id="docs"></a> Automatic documentation

Additionally, you can make simple documentation for your library. Use the following format to add documentation to your functions:

```go
/*
    Returns a greeting to the name given.
    func sayHello(name string) string
*/
func sayHello(name string) (val interface{}, err error) {
    return "Hello, " + name, err
}
```

This will automatically be added to a markdown file in the same directory by running the `autodocs.py` file. It will result in this:

> ## **`sayHello`**
>
> Returns a greeting to the name given.
>
> ```go
> func sayHello(name string) string
> ```
