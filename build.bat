@echo off
if not exist ./bin (
    mkdir bin
)

go build -o ./bin/fizz.exe main.go
echo done