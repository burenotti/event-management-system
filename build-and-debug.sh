#!/bin/ash

go build -gcflags="all=-N -l" -o /tmp/main ./cmd/app/main.go
dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./tmp/main
