#!/bin/bash
# start-dev.sh

# Kill existing processes
killall hugo 2>/dev/null
pkill -f "go run main.go" 2>/dev/null

# Start Hugo
cd ~/blogg
hugo server --port 1313 --bind 127.0.0.1 &

# Wait for Hugo to start
sleep 2

# Start Go server
cd api
go run main.go