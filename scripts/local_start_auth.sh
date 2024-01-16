#! /bin/bash

# Array of directories where your Go modules are located
declare -a goModuleDirs=("services/api/register" "services/api/login")

# Loop through each directory
for dir in "${goModuleDirs[@]}"; do
    echo "Starting auth service in $dir"
    (
        cd "$dir" || exit                 # Change to the directory, exit if it fails
        go run main.go 2>&1 >"logs.log" & # Run the Go server and redirect output to logs.log
    )
done

echo "All servers started."
