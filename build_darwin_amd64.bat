mkdir build
mkdir "build/darwin_amd64"
set GOARCH=amd64
set GOOS=darwin
go build -o build/darwin_amd64/tender -ldflags "-s -w" cmd/tender/main.go
pause