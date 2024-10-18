mkdir build
mkdir "build/windows_386"
set GOARCH=386
set GOOS=windows
go build -o build/windows_386/tender.exe  -ldflags "-s -w" cmd/tender/main.go
pause