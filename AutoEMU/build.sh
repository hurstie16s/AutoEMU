#!/bin/bash

ARG1="${1}"

if [ "${ARG1}" = "config" ]; then
    sudo rm -f config.yaml
fi

# Remove previous build
sudo rm -f AutoEMU
# Remove previous log file
sudo rm -f Logs/*.log
# Remove terminal tests
sudo rm -r WorkSpace/FirmwareFiles/terminalTest

# Build go project
go build .
