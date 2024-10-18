mkdir build
mkdir "build/linux_amd64"
set GOARCH=amd64
set GOOS=linux
go build -o build/linux_amd64/tender -ldflags "-s -w" cmd/tender/main.go
pause