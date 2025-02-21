#!bin/bash

# find the port id of the running container
PID=`pgrep authApp`
# attach delve to the running container
dlv attach $PID --headless --listen=:2345 --api-version=2 --log