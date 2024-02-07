#!/bin/bash

# Check if two arguments are provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 start_port end_port"
    exit 1
fi

start_port=$1
end_port=$2

# Loop through the range of ports
for port in $(seq "$start_port" "$end_port"); do
    # Use lsof to find processes listening on the port
    # Note: This requires root privileges to see all processes
    pid=$(lsof -ti tcp:"$port")

    # If a PID is found, kill the process
    if [ ! -z "$pid" ]; then
        echo "Killing process on port $port with PID $pid"
        kill -SIGTERM "$pid" || echo "Failed to kill process $pid on port $port"
    fi
done
