#! /bin/bash

# Array of directories where your Go modules are located
declare -a goModuleDirs=("services/api/getProperty" "services/api/getProperties" "services/api/addProperty")

# Loop through each directory
for dir in "${goModuleDirs[@]}"; do
    echo "Starting service in $dir"
    (
        cd "$dir" || exit # Change to the directory, exit if it fails
        service_name="/" read -ra ADDR <<<"$dir"
        go run main.go 2>&1 >"${ADDR[-1]}.logs.log" & # Run the Go server and redirect output to logs.log
    )
done

echo "All servers started."
