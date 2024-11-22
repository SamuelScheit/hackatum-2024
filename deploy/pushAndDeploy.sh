#!/bin/bash
cd /Users/jgs/Developer/Projekte/hackatum/hackatum24-check24/
git add .
git commit -m "push and deploy"
git push origin 
./deploy/deploy.sh
