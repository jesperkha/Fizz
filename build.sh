set -e

# Include all libraries
python build_lib.py
# Lint
go vet .
# Create bin folder if it doesnt exist and build
[ ! -d "./bin" ] && mkdir bin
go build -o ./bin/fizz.exe run.go main.go