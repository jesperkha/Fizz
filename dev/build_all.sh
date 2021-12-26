set -e
[ ! -d "./bin" ] && mkdir bin
python build_lib.py
go vet .
GOOS=windows GOARCH=amd64 go build -o ./bin/fizz.exe .
GOOS=linux GOARCH=amd64 go build -o ./bin/fizz .
