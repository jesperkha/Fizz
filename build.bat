@echo off
if not exist ./bin (
    mkdir bin
)

go build -o ./bin/fizz.exe run.go interp.go main.go