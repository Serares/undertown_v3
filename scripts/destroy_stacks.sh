#! /bin/bash

# Destroy the distributions first
# Destoru the raw-images-bucket
# Destory the processed-images-bucket
# All the other stacks then

# Define your stack names in order of deletion, starting with the most dependent
STACKS=(
    "SharedOriginsStack-dev"
    "AdminDistribution-dev"
    "HomeDistribution-dev"
    "raw-images-bucket-dev"
    "processed-images-bucket-dev"
)

for stack in "${STACKS[@]}"; do
    echo "Attempting to delete stack: $stack"
    aws cloudformation delete-stack --stack-name $stack

    echo "Waiting for stack to be deleted..."
    aws cloudformation wait stack-delete-complete --stack-name $stack
    if [ $? -eq 0 ]; then
        echo "Stack $stack deleted successfully."
    else
        echo "Error deleting stack $stack."
        exit 1
    fi
done

# Use the cdk to delete all the remaining stacks

cd cdk && cdk destroy --all --require-approval never -y

echo "All stacks deleted successfully."
