# <a id="top"></a> Fizz language documentation

**Language specifics:**

- [Overview](#overview)
- [Grammar](#grammar)
- [Types](#types)
- [Keywords](#keywords)
- [Expressions and operators](#expressions-and-operators)

**Variables and printing:**

- [Print and Type](#print-and-type)
- [Error and Exit](#error-and-exit)
- [Variables](#variables)
- [Enums](#enums)

**Conditionals and loops:**

- [If statements and logic](#if-statements-and-logic)
- [While loop](#while-loop)
- [Repeat loop](#repeat-loop)
- [Range loop](#range-loop)
- [Break and skip](#break-and-skip)

**Objects:**

- [Functions](#functions)
- [Objects](#objects)
- [Arrays](#arrays)
- [Reference](#reference)

**Files and imports**

- [File imports](#file-imports)
- [Libraries](#libraries)

<br>

> **IMPORTANT:**
> In the examples below, `//` is used for comments since Go is used as the language for the codeblocks (for better highlighting). However, in Fizz, `#` signifies a comment.

<br>

## Overview

- Fizz is dynamically typed, meaning type checks are only performed at runtime
- Comments are created with a hashtag `#` (as mentioned above) and end at the first found newline
- All statements that do not have a block tied to them must end with a semicolon `;`. This, of course, means you can have all your code on one line

<br>

## Grammar

[A definitive grammar sheet](./grammar.md)

<br>

## Types

Fizz is strongly typed, meaning unmatched types in certain expressions will cause a runtime error.

- `string` Any string of text with a starting and ending quote `"` symbol. Can span over multiple lines. Can also include `\n` for a new line, or `\t` for a tab.

- `number` Any number, including floats.

- `nil` No value.

- `bool` Keywords `true` and `false`.

  - **Truthyness**: Any expression that does not evaluate to `nil` or `false` is truthy.

- `object` Type of object instance

- `function` Type of function or object constructor

- `array` Type of array instance

<br>

## Keywords

Keyword names are reserved and cannot be used for variable names. Here is a list of all keywords in Fizz:

```
var       print     type       func      range     error
exit      skip      break      return    in
false     nil       include    if        enum
import    define    true       while     repeat
```

<br>

## Expressions and operators

Fizz features a lot of standard syntax similar to other languages. For example, all normal expressions using the basic arithmetic and logic operators will work in Fizz, including the modulo operator and the hat operator. Plus can also be used for joining strings.

Operators:

- Binary operators:
  ```go
  +   -   *   /   %   ^   <
  >   ==  !=  >=  <=  &   :  in
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

## Print and Type

In Fizz, `print` is a _statement_, not a function. However, `type` is an _operator_, not a function, and gives a string value.

```go
print "Hello";      // Hello
print type "World"; // string
```

<br>

## Error and Exit

The `error` statement prints out a message (or value) as an error and exits the program.

```go
error "some error occured"; // prints message as error and exits
```

Theres also an `exit` statement. This will just print the value out (same as `print`) and then exit with no error. If an expression is not given, `exit` will just quit without printing anything.

```go
exit "goodbye"; // prints message and exits
```

<br>

## Variables

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

## Enums

Fizz also has enums to quickly create a lot of variables with unique values.

```go
enum {
  banana
  apple
  orange
}

enum {
  pear
}

print banana; // 0
print apple;  // 1
print orange; // 2
print pear;   // 0
```

<br>

## If statements and logic

Fizz features simple `if` and `else` statements, but not `else if`. The 'and' operator is `&` and 'or' is `:`.

```go
height := 172;

if height > 158.8 {
    print "Taller than Kevin Hart";
} else {
    print "Not taller than Kevin Hart";
}
```

<br>

## While loop

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

## Repeat loop

A repeat loop is another basic flow controller that just executes a block n times.

```go
// prints "hi" 5 times
repeat 5 {
  print "hi";
}
```

<br>

## Range loop

The range loop is a more advanced form of loop, kind of a hybrid bewteen pythons `for _ in _` statements and other languages `for` loops. The simplest use case is to just give one argument to the right side of `in`.

In this case, it will just loop with `n` going from 0 to 9, as the default starting number is 0:

```go
range n in 10 {
  print n;
}
```

Providing two arguments defines both the start and end for the loop:

```go
// Goes from 3 to 7
range n in 3, 8 {
  print n;
}
```

Three arguments define start, end, and iteration amount. The amount can be negative and make the loop count down, but if the conditions are set in a way where the loop will never end, an error is raised:

```go
range n in 0, 5, 0.5 {
  print n;
}

// error: infinite loop not allowed for range statement
range q in 0, 5, -1 {

}
```

Additionally, you can range over an array:

```go
days := ["Monday", "Tuesday", "Wednesday"];

range day in days {
  print day;
}
```

<br>

## Break and skip

- `skip` skips to next iteration in loop
- `break` breaks out of loop

<br>

## Functions

You can declare a function using the `func` keyword. Functions will return `nil` if no other return value is specified. Passing an incorrect argument number will cause a runtime error.

```go
func add(a, b) {
    return a + b;
}

print add(5, 2); // 7
```

<br>

## Objects

Object structures can be defined with the `define` keyword. This creates an object template which you can use to make your own structured data. The field names can be separated by a line break, comma, or space. The fields of the object do not have a specific type, unlike languages like C and Go. Object values support reassignment too.

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

Under the hood, the `define` statement creates a function that returns an object with the specified values. That means `Person` is a `function` type and `john` is an `object` type:

```go
print type Person // function
print type john   // object
```

<br>

## Arrays

Arrays in Fizz are just an array of values, of which can be any type. You get the value of a specific index in an array by using the index getter syntax. Indexes begin at 0. Additionally, you can get the length of an array with the built-in `len` function.

```go
names := ["John", "Susan", "Carl"];
print names[0]; // John

names[2] = "Timmy";
print names; // ["John", "Susan", "Timmy"]

print len(names); // 3
```

You can use the `in` operator to check if an element is present in an array:

```js
print "dog" in ["cat", "dog", "fox"]; // true
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

## Reference

In Fizz, objects and arrays are passed by reference. This means you can modify them directly when passing them as a function argument:

```go
food := ["bread", "pasta", "rice"];

func addCarrot(arr) {
  push(arr, "carrot");
}

addCarrot(food);
print food; // ["bread", "pasta", "rice", "carrot"]
```

```go
define Phone {
  brand
  version
}

iphone := Phone("Apple", 1);

func upgrade(phone) {
  phone.version += 1;
}

upgrade(iphone);
print iphone.version; // 2
```

<br>

## File imports

You can import files by using the `import` statement. The given path, or name, is always relative to the file that the program started in. Circular imports are not allowed and an error will be raised if one is found. Importing creates an object with all the values of the imported file. The object is declared with the name of the file that was imported, so files with the same names cannot be imported in the same file. (in the future `import x as y` syntax will be added to fix this)

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

## Libraries

Fizz libraries are different from imports. They are not other Fizz files, but rather Go files. This is to make it possible for functionality to be added to Fizz without baking it straight in. You can read the [library documentation](./libraries.md) to find out how they work and how to create your own.

Fizz has a standard library built in. Include it with the `include` keyword. The functions are documented in the `lib/<module_name>` directory.

```go
include "str";

age := 32;
print "John is " + str.toString(age); // prints: John is 32
// Causes no type error because age is converted to a string
```

```go
include "io";
include "str";

// Will prompt the terminal and wait for user input.
// Continues at newline or exit with ctrl-C.
meters := io.input("Enter height in meters: ");

feet := str.toNumber(meters) * 3.281;
print "You are: " + str.toString(feet) + " feet tall";
// Built-in string formatting is on the todo list, don't worry ;)
```

<br>

[Go to top](#top)
