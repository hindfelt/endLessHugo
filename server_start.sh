#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}Starting development environment...${NC}"

# Kill existing processes
echo -e "${GREEN}Killing existing processes...${NC}"
killall hugo 2>/dev/null
killall caddy 2>/dev/null
pkill -f "go run main.go" 2>/dev/null

# Wait for ports to be freed
sleep 2

# Start Hugo
echo -e "${GREEN}Starting Hugo...${NC}"
cd ~/blogg  # Adjust this path
hugo server --port 1313 --bind 127.0.0.1 > /dev/null 2>&1 &

# Wait for Hugo to start
sleep 2

# Start Go auth server
echo -e "${GREEN}Starting Go auth server...${NC}"
cd ~/blogg/api  # Adjust this path
go run main.go > /dev/null 2>&1 &

# Wait for Go server to start
sleep 2

# Start Caddy
echo -e "${GREEN}Starting Caddy...${NC}"
caddy run --config Caddyfile > /dev/null 2>&1 &

# Wait for Caddy to start
sleep 2

# Check if services are running
echo -e "${GREEN}Checking services...${NC}"
if nc -z localhost 1313; then
    echo "✅ Hugo is running on port 1313"
else
    echo "❌ Hugo failed to start"
fi

if nc -z localhost 8080; then
    echo "✅ Go auth server is running on port 8080"
else
    echo "❌ Go auth server failed to start"
fi

if nc -z localhost 8000; then
    echo "✅ Caddy is running on port 8000"
else
    echo "❌ Caddy failed to start"
fi

echo -e "${GREEN}Development environment started!${NC}"
echo "📝 Access your site at: http://localhost:3000"

# Keep script running and show logs
echo "Press Ctrl+C to stop all services"
wait