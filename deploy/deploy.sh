#!/bin/bash
ssh -J jumper@jumphost.hackatum.check24.de challenger@team13.hackatum.check24.de "cd ./hackatum-2024/server; git pull; chmod +x ./deploy/buildAndRestart.sh; sudo ./buildAndRestart.sh"