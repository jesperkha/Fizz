# **Fizz**

- [Overview](#overview)
- [Expressions](#expr)
- [Print and Type](#prt)
- [Variables](#var)
- [If statements and logic](#if)
- [While and repeat](#while)
- [Break, skip, and exit](#break)
- [Functions](#func)
- [Objects](#obj)

<br>

> **IMPORTANT:**
> In the examples below, `//` is used for comments since Go is used as the language for the codeblocks (for better highlighting). However, in Fizz, `#` signifies a comment.

<br>

### <a id="overview"></a> **Language overview**

Heres a list of things about Fizz that you should know before reading the rest of this page or trying Fizz out for yourself:

- Fizz is dynamically typed, meaning type checks are only performed for expression evaluation
- Comments are created with a hashtag `#` (as mentioned above) and end at the first found line break
- All statements that do not have a block tied to it must end with a semicolon `;`. This, off course, means you can have all your code on one line

<br>

### <a id="expr"></a> **Expressions**

Fizz features a lot of standard syntax similar to other languages. For example, all normal expressions using the basic arithmetic and logic operators will work in Fizz, including the modulo operator and the hat operator. Plus can also be used for joining strings.

```go
4 % 2 == 0;
(3 ^ 2) == 9;
"Hello" + "World";
```

Fizz is dynamically typed, but will not convert types in expressions. Instead and error is raised when types do not match.

<br>

### <a id="prt"></a> **Print and Type**

In Fizz, `print` is a _statement_, not a function. However, `type` is an _operator_, not a function, and gives a string value.

```go
print "Hello";
print type "World";
```

<br>

### <a id="var"></a> **Variables**

You can declare a variable using the `var` statement. The value can be re-assigned later and even change type.

```go
var name = "John";
name = "Carl";
name = 3;

// Error, 'name' is already defined
var name = "Susan";
```

Local variables override higher level scope:

```go
var age = 10;

{
	// Overrides flobal 'age' variable
	var age = 20;
}
```

You can use shorthand assignment operators too:

```go
var n = 1;
n += 2;
n -= 2;
n *= 2;
n /= 2;
```

You can also use the `+=` operator with strings.

<br>

### <a id="if"></a> **If statements and logic**

Fizz features simple if and else statements, but not else-if. The 'and' operator is `&` and 'or' is `:`.

```go
var height = 172;

if age > 158.8 {
    print "Taller than Kevin Hart";
} else {
    print "Not taller than Kevin Hart";
}
```

<br>

### <a id="while"></a> **While and repeat**

Fizz has a while statement similar to most other languages. If you leave the expression field empty it will just run forever.

```js

while n < 10 {
    // loops until condition is false
}

while {
    // loops until break or program exit
}
```

(Temporary: will be replaced with a range statement) The repeat statement is a little different. It is a condenced 'for' loop. You first declare a variable name followed by a legal repeat operator and range. Currently, only `<` is allowed.

```go
repeat n < 10 {
    // (creates n) loops 10 times, incrementing n
}
```

<br>

### <a id="break"></a> **Break, skip, and exit**

- `skip` skips to next iteration in loop
- `break` breaks out of loop
- `exit` stops program execution

<br>

### <a id="func"></a> **Functions**

You can declare a function using the `func` keyword.

```go
func add(a, b) {
    return a + b;
}

print add(5, 2); // 7
```

<br>

### <a id="obj"></a> **Objects**

Object structures can be defined with the `define` keyword. This creates a object template which you can use to make your own structured data. The fields of the object do not have a specific type, unlike languages like C and Go. Object values support reassignment too.

```js
define Person {
    name
    age
}

var john = Person("John", 31);
print john.name; // "John"

john.age = 99;
print john.age; // 99
```
