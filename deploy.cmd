@echo off
cls
rmdir /q /s dest
mkdir dest

copy default.config.json dest\config.json
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64
go build -o dest/ping.exe
cd dest
"c:\program files\7-zip\7z.exe" a ping_windows_x64.zip ping.exe config.json
cd ..

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o dest/ping
cd dest
"c:\program files\7-zip\7z.exe" a ping_linux_amd64.zip ping config.json
cd ..

del dest\config.json
del dest\ping
del dest\ping.exe



