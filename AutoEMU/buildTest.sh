#!/bin/bash

sudo rm -f AutoEMU
sudo rm -f Logs/*.log
sudo rm -r WorkSpace/FirmwareFiles/unitTestMachine

go test
