#!/bin/ash

go test -v -count=1 -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o ./tmp/coverage.html
rm coverage.out
