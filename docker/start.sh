#!/usr/bin/env sh

# start processes
themisd start > ./logs/themisd.log &
themisd rest-server > ./logs/themisd-rest-server.log &
sleep 100
bridge start --all > ./logs/bridge.log &

# tail logs
tail -f ./logs/themisd.log ./logs/themisd-rest-server.log ./logs/bridge.log
