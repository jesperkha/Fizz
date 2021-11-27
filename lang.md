### **Keywords**

- `var`
- `func`
- `print`
- `type`
- `repeat`
- `while`
- `return`
- `break`
- `skip`
- `exit`

<br>

### **Types**

- `float64`
- `bool`
- `string`
- `function`
- `nil`

<br>

### **Syntax**

Fizz features a lot of standard syntax similar to other languages. For example, all normal expressions using the basic arithmatic and logic operators will work in Fizz:

`1 + 1 == 2` <br>
`(15 / 3) > 4`

For conditions and flow control there are some things to take note of. Firstly, you can create and infinite while loop by not giving it an expression:

```
while {
    # Will run forever
}
```

Some other unique traits include the `repeat` statement:

```
repeat n < 10 {
    # executes block 10 times
    print n; # 0 -> 9
}
```

This is really just a simplified for loop. As of now it only works with the less than operator. In the future more will be added. Below is a list of code examples to further explain the syntax of Fizz. They should be very straight forward and easy to understand for anyone experienced with similar programming languages.

<br>

### **Code examples**

Some code examples that show most, if not all, of the syntax Fizz features:

```
func add(a, b) {
    return a + b;
}

var num1 = 3;
var num2 = 4;
print add(num1, num2); # 7
```

```
func isEven(num) {
    return num % 2 == 0;
}

repeat n < 10 {
    if isEven(n) {
        print "Even!";
    } else {
        print "Odd!";
    }
}
```

```
var name = "John";
var age = 31;

print type name == type age; # false
```

```
var n = 0;
while n < 5 {
    n += 1;
}

while {
    if n == 0 {
        break;
    }

    n -= 1;
}
```

```
var n = 1;
var done = false;

while !done {
    n *= 2;
    if n < 100 {
        skip;
    }

    exit;
}

print "This will not be printed";
```

```
var name = "Bob";
var age = 46;
var job = "Pilot";

if name == "Bob" & job == "Pilot" {
    print "Bob is a pilot";
}

if age == 42 : age == 46 {
    print "Either 42 or 46";
}

```
