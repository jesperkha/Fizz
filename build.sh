set -e

echo -ne "Adding libraries\t"
python lib/build.py
echo "done"

echo -ne "Generating docs \t"
python lib/autodocs.py
echo "done"

echo -ne "Linting code    \t"
go vet .
echo "done"

echo -ne "Building binary \t"
[ ! -d "./bin" ] && mkdir bin
go build -o ./bin/fizz.exe .
echo "done"