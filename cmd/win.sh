#!/bin/bash

# shellcheck disable=SC2034

export GOOS=windows
go build .
./cashflow-go.exe
