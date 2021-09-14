#!/bin/bash
rm -rf docs
go mod init
swag init -g communication/http/http.go
go mod tidy
go build
