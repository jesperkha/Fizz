# Fizz update changelog

## Version 1.1.0

**Whats new:**

- Newline and tab characters

**Bug fixes and more:**

## Version 1.0.0

**Whats new:**

- Enums
- New repeat statement
- Range statement

**Bug fixes and more:**

- Changed env to allow acces to definitions / reassignments _after_ a closure was formed
- Added new recursive equality check for objects and arrays
- Implemented callstack and -f flag to show it
- Error for exceeding recursion limit
- Fixed bug where environments would be referenced and not copied

## Version 0.6.0

**Whats new:**

- Libraries
- Error statement

**Bug fixes and more:**

- Automatic documentation for libraries

<br>

## Version 0.5.0

**Whats new:**

- Arrays
- New `:=` operator for variable declaration, removed `var` statement (temporary)

**Bug fixes and more:**

- Prettier print for values
- Fixed the semicolon error to now actually show when there is a semicolon missing, instead of just giving an expression error
- Fixed error that would occur when calling group expressions

<br>

## Version 0.4.0

**Whats new:**

- Closures

**Bug fixes and more:**

- Errors for circular imports
- Patched bug where env was cleared in terminal mode

<br>

## Version 0.3.0

**Whats new:**

- Added file imports

**Bug fixes and more:**

- Fixed error messages for break, skip, and return errors
- Function error traceback to origin file
- New env structure
- Added changelog
