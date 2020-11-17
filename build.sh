#! /bin/bash

#macOS版本
go build -o bin/macos/magpie ./main.go
go build -o bin/macos/magctl ./magctl/main.go

#Linux版本
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/magpie ./main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/magctl ./magctl/main.go

#Windows版本
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows/magpie.exe ./main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows/magctl.exe ./magctl/main.go