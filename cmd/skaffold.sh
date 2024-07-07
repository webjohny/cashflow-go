#!/bin/bash
# skaffold.sh

kubectl config use-context do-fra1-cashflow-k8s-1-30-1-do-0-fra1-1716791288523
skaffold dev -vdebug --default-repo=registry.digitalocean.com/cashflow-registry