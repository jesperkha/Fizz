set -e
python build_lib.py
go vet .
[ ! -d "./bin" ] && mkdir bin
go build -o ./bin/fizz.exe run.go main.go