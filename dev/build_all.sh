set -e
[ ! -d "./bin" ] && mkdir bin
python lib/build.py
python lib/autodocs.py
go vet .
GOOS=windows GOARCH=amd64 go build -o ./bin/fizz.exe .
GOOS=linux GOARCH=amd64 go build -o ./bin/fizz .
