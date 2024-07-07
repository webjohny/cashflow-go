#!/bin/bash

# Check if deploying to production
if [[ "$KUBE_ENV" == "production" ]]; then
  echo "Error: Deploying to production is not allowed through this script."
  exit 1
fi

# Check for risky environment variables
if [[ "$DB_DROP_ALL" == "true" ]]; then
  echo "Error: DB_DROP_ALL is set to true. Aborting deployment to prevent data loss."
  exit 1
fi

echo "Pre-deploy checks passed."