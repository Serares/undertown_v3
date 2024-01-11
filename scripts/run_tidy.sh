#! /bin/bash

# Array of directories where your Go modules are located
declare -a goModuleDirs=("services/api/*/" "services/ssr/*/")

for path in "${goModuleDirs[@]}"; do
    # Loop through each directory
    for dir in $path; do
        echo "Running tidy in dir: $dir \n"
        (
            cd "$dir" || exit # Change to the directory, exit if it fails
            go mod tidy
        )
    done
done
echo "Tidy done"