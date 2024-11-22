#!/bin/bash
ssh -J jumper@jumphost.hackatum.check24.de challenger@team13.hackatum.check24.de "cd ./hackatum-2024; git pull; chmod +x ./server/service.sh; sudo ./server/service.sh"