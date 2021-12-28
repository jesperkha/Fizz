set -e

echo adding libraries...
python lib/build.py

echo linting code...
go vet .

echo building binary...
[ ! -d "./bin" ] && mkdir bin
go build -o ./bin/fizz.exe run.go main.go

echo "done!"