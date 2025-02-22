#!bin/bash

# find the port id of the running container
PID=`pgrep brokerApp`
# attach delve to the running container
dlv attach $PID --headless --listen=:2345 --accept-multiclient --api-version=2 --log