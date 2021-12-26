set -e

echo adding libraries...
python build_lib.py

echo linting...
go vet .

echo building binary...
[ ! -d "./bin" ] && mkdir bin
go build -o ./bin/fizz.exe run.go main.go