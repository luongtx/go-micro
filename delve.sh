#!/bin/bash
# install go and delve
apk update && apk add go musl-dev
go install github.com/go-delve/delve/cmd/dlv@latest
export PATH=$PATH:$HOME/go/bin
# find the port id of the running container
TARGET_BINARY=$1
PID=`pgrep $TARGET_BINARY`
# attach delve to the running container
dlv attach $PID --headless --listen=:2345 --accept-multiclient --api-version=2 --log