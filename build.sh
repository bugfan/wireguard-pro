#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wireguard-go
scp wireguard-go root@app.lt53.cn:/home