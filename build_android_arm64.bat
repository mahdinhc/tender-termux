mkdir build
mkdir "build/android_arm64"
set GOARCH=arm64
set GOOS=android
go build -o build/android_arm64/tender -ldflags "-s -w" cli/tender/main.go
pause