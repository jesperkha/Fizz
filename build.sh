set -e

echo -n "adding libraries... "
python lib/build.py
echo "done"

echo -n "generating docs... "
python lib/autodocs.py
echo "done"

echo -n "linting code... "
go vet .
echo "done"

echo -n "building binary... "
[ ! -d "./bin" ] && mkdir bin
go build -o ./bin/fizz.exe .
echo "done"