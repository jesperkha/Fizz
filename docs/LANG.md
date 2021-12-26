# <a id="top"></a> **Fizz language documentation**

- [Overview](#overview)
- [Grammar](#grammar)
- [Types](#types)
- [Keywords](#keywords)
- [Expressions and operators](#expr)
- [Print and Type](#prt)
- [Variables](#var)
- [If statements and logic](#if)
- [While loop](#while)
- [Break and skip](#break)
- [Functions](#func)
- [Objects](#obj)
- [Arrays](#arrays)
- [File imports](#import)
- [Libraries](#libraries)

<br>

> **IMPORTANT:**
> In the examples below, `//` is used for comments since Go is used as the language for the codeblocks (for better highlighting). However, in Fizz, `#` signifies a comment.

<br>

## <a id="overview"></a> **Language overview**

- Fizz is dynamically typed, meaning type checks are only performed at runtime
- Comments are created with a hashtag `#` (as mentioned above) and end at the first found newline
- All statements that do not have a block tied to them must end with a semicolon `;`. This, of course, means you can have all your code on one line

<br>

## <a id="grammar"></a> **Grammar**

A definitive grammar sheet can be found [here](./grammar.md).

<br>

## <a id="types"></a> **Types**

Fizz is strongly typed, meaning unmatched types in certain expressions will cause a runtime error.

- `string` Any string of text with a starting and ending quote `"` symbol.

- `number` Any number, including floats.

- `nil` No value.

- `bool` Keywords `true` and `false`.

  - **Truthyness**: Any expression that does not evaluate to `nil` or `false` is truthy.

- `object` Type of object instance

- `function` Type of function or object constructor

<br>

## <a id="keywords"></a> **Keywords**

Keyword names are reserved and cannot be used for variable names. Here is a list of all keywords in Fizz:

```
var       print     type       func
exit      skip      break      return
false     nil       include    if
import    define    true       while
```

<br>

## <a id="expr"></a> **Expressions and operators**

Fizz features a lot of standard syntax similar to other languages. For example, all normal expressions using the basic arithmetic and logic operators will work in Fizz, including the modulo operator and the hat operator. Plus can also be used for joining strings.

Operators:

- Binary operators:
  ```go
  +   -   *   /   %   ^   <
  >   ==  !=  >=  <=  &   :
  ```
- Unary operators:
  ```go
  -   !  type
  ```
- Assignment operators:
  ```go
  =   :=
  ```

<br>

## <a id="prt"></a> **Print and Type**

In Fizz, `print` is a _statement_, not a function. However, `type` is an _operator_, not a function, and gives a string value.

```go
print "Hello";
print type "World";
```

Theres also an `exit` statement. It's almost identical to `print`, but it also exits the program at execution. If an expression is not given, `exit` will just quit without printing anything.

```js
exit "goodbye"; // prints 'goodbye' and exits program
```

<br>

## <a id="var"></a> **Variables**

You can declare a variable using the `:=` operator. The value can be re-assigned later and even change type.

```go
name := "John";
name = "Carl";
name = 3;

// Error, 'name' is already defined
name := "Susan";
```

Local variables override higher level scopes:

```go
age := 10;

{
	// Overrides global 'age' variable
	age := 20;
}
```

You can use shorthand assignment operators too:

```go
n := 1;
n += 2;
n -= 2;
n *= 2;
n /= 2;
```

You can also use the `+=` operator with strings.

<br>

## <a id="if"></a> **If statements and logic**

Fizz features simple if and else statements, but not else-if. The 'and' operator is `&` and 'or' is `:`.

```go
height := 172;

if age > 158.8 {
    print "Taller than Kevin Hart";
} else {
    print "Not taller than Kevin Hart";
}
```

<br>

## <a id="while"></a> **While loop**

Fizz has a while statement similar to most other languages. If you leave the expression field empty it will just run forever.

```js
while n < 10 {
    // loops until condition is false
}

while {
    // loops until break or program exit
}
```

<br>

## <a id="break"></a> **Break and skip**

- `skip` skips to next iteration in loop
- `break` breaks out of loop

<br>

## <a id="func"></a> **Functions**

You can declare a function using the `func` keyword. Functions will return `nil` if no other return value is specified. Passing an incorrect argument number will cause a runtime error.

```go
func add(a, b) {
    return a + b;
}

print add(5, 2); // 7
```

<br>

## <a id="obj"></a> **Objects**

Object structures can be defined with the `define` keyword. This creates a object template which you can use to make your own structured data. The fields of the object do not have a specific type, unlike languages like C and Go. Object values support reassignment too.

```go
define Person {
    name
    age
}

john := Person("John", 31);
print john.name; // "John"

john.age = 99;
print john.age; // 99
```

Under the hood, the `define` statement creates a function that returns an object with the specified values. That means `Person` is a function type and `john` is an object type.

<br>

## <a id="arrays"></a> **Arrays**

Arrays in Fizz are just an array of values, of which can be any type. You get get the value of a specific index in an array by using the index getter syntax. Indexes begin at 0. Additionally, you can get the length of an array with the built in `len` function.

```go
names := ["John", "Susan", "Carl"];
print names[0]; // John

names[2] = "Timmy";
print names; // ["John", "Susan", "Timmy"]

print len(name); // 3
```

### Push

To push elements into an array use the built-in `push` function:

```go
arr := [1, 2, 3];
push(arr, 4);

print arr; // [1, 2, 3, 4]
```

### Pop

There is also a `pop` function to remove and return the last element:

```go
arr := [1, 2, 3];
print pop(arr); // 3
print arr;      // [1, 2]
```

<br>

## <a id="import"></a> **File imports**

You can import files by using the `import` statement. The given path, or name, is always relative to the file that the program started in. Circular imports are not allowed and an error will be raised if one is found. The imported object name is the filename, so files with the same names cannot be imported in the same file. (in the future `import x as y` syntax will be added to fix this)

```go
// other.fizz

name := "John";
```

```js
// main.fizz

import "other";
// Also valid:
// import "other.fizz";
// import "./other.fizz";

print other.name;
```

```console
$ ls
main.fizz   other.fizz

$ fizz main
John
```

<br>

## <a id="libraries"></a> **Libraries**

<!-- Todo: remove in v6 -->

> Version 0.6.0 or higher

Fizz libraries are different from imports. They are not other Fizz files, but rather Go files. This is to make it possible for functionality to be added to Fizz without baking it straight in. You can read the [library documentation](./libraries.md) to find out how they work and how to create your own.

Fizz has a standard library built in. To use it, use the `include` keyword.

```go
include "std";

age := 32;
print "John is " + std.toString(age); // prints: John is 32
// Causes no type error because age is converted to a string
```

```go
include "std";

// Will prompt the terminal and wait for user input.
// Continues at newline or exit with ctrl-C.
meters := std.input("Enter height in meters: ");

feet := std.toNumber(meters) * 3.281;
print "You are: " + std.toString(feet) + " feet tall";
// Built in string formatting will come soon dont worry ;)
```

<br>

[Go to top](#top)
