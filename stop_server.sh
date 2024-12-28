#!/bin/bash
killall hugo 2>/dev/null
killall caddy 2>/dev/null
pkill -f "go run main.go" 2>/dev/null
echo "Development environment stopped"