set -e
gofmt -w -s -d .
go vet .
go mod tidy
go mod verify
git add .
git commit -m $1