mkdir build
mkdir "build/windows_amd64"
set GOARCH=amd64
set GOOS=windows
go build -o build/windows_amd64/tender.exe  -ldflags "-s -w" cli/main.go
pause