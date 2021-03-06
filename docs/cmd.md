# Fizz command line use

You run a Fizz file by running the `fizz` program followed by a path to a file or just the filename if it is in the current directory. Note that you don't need to include the `.fizz` suffix, it's optional.

```console
$ fizz myFile
```

As also mentioned in the readme, running `fizz` with no arguments runs the terminal mode. You can then write any valid Fizz code and run it live. Errors do not terminate the session.

```console
$ fizz
type 'exit' to terminate session
1 : print "hello";
hello
```

Creating a new line after an opening brace `{` will auto indent for you and not execute the code until you close it again with a closing brace `}`:

```console
$ fizz
type 'exit' to terminate session
1 : func main() {
2 :     print "hello";
3 : }
4 :
5 : main();
hello
```

You can at any point type `exit` followed by enter to close the program. Using `ctrl-C` is also possible, but not recommended.

<br>

## Flags

There are multiple flags you can use, however, some will only take effect when running a file.

<br>

**Info flags**

- `--help` print information on how to use the program and also all available flags
- `--version` print the version of the program

<br>

**Config flags**

- `-f` print function callstack upon error
- `-e` print the global environment after program finish

<br>

## Subcommands

- **`help`**

  ```console
  $ fizz help [command]
  ```

  Prints out a help message for the given command.

- **`docs`**

  ```console
  $ fizz docs [lib name]
  ```

  When given a name of a valid library (`str`, `io` etc), it will print out all of the functions defined in that library. If you have made your own library and some functions are not showing up make sure to re-run `autodocs.py`.
