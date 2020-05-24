#!/bin/bash
if [ ! -f secrets/app.ini ]; then
   cp conf/app.ini secrets/
fi
go run main.go
