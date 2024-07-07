#!/bin/bash
# entrypoint.sh

# Create or clear the .env file
echo "" > /app/.env

# Read each file in the directory, assuming they are from the mounted secrets
for filename in /etc/secrets/*; do
    # Extract the content of each file
    content=$(cat "$filename")
    filename=$(basename "$filename")

    echo "${filename}=${content}" >> /app/.env
done

# Execute the main command, e.g., start your application
exec "$@"