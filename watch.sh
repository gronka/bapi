#!/bin/sh

export GO111MODULE=on  # must activate this for modules in gopath
# go mod init

while true; do
  #go run *.go
  go build
  $@ &
  PID=$!
  inotifywait -r -e modify .
  kill $PID
done
