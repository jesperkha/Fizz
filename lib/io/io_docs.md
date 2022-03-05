<!--
func input(prompt string) string
func readFile(filename string) string
func writeFile(filename string, content string)
func appendFile(filename string, content string)
func readDir(dir string) []string
func curDir() string
func exists(filename string) bool
func newDir(name string)
func newFile(name string)
-->

# Methods in io library

## **`input`**

Gets user input from stdin.

```go
func input(prompt string) string
```

<br>

## **`readFile`**

Reads file and returns text.

```go
func readFile(filename string) string
```

<br>

## **`writeFile`**

Writes content to file. Overwrites previous file content.

```go
func writeFile(filename string, content string)
```

<br>

## **`appendFile`**

Appends content to file.

```go
func appendFile(filename string, content string)
```

<br>

## **`readDir`**

Returns list of files/directories in dir.

```go
func readDir(dir string) []string
```

<br>

## **`curDir`**

Returns current working directory

```go
func curDir() string
```

<br>

## **`exists`**

Returns true if file exists.

```go
func exists(filename string) bool
```

<br>

## **`newDir`**

Creates new directory.

```go
func newDir(name string)
```

<br>

## **`newFile`**

Creates new file. If the file already exists it will be overwritten.

```go
func newFile(name string)
```

<br>

