#!bin/bash
# install go and delve
apk update && apk add go musl-dev
go install github.com/go-delve/delve/cmd/dlv@latest
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$PATH
# find the port id of the running container
PID=`pgrep authApp`
# attach delve to the running container
dlv attach $PID --headless --listen=:2345 --accept-multiclient --api-version=2 --log