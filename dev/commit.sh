set -e
go test ./test
go mod tidy
go mod verify
go vet .
gofmt -w -s .
python lib/autodocs.py
git add .
git commit -m "$1"