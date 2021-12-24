set -e
[ ! -d "./bin" ] && mkdir bin
python build_lib.py
go vet .
GOOS=windows GOARCH=amd64 go build -o ./export/fizz.exe main.go run.go
GOOS=linux GOARCH=amd64 go build -o ./export/fizz main.go run.go