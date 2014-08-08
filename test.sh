#!/bin/bash

echo "Killing off etcd if it is running..."
killall etcd

echo "Starting etcd..."
nohup etcd > /dev/null &
sleep 3

echo "Adding test variables..."
etcdctl set /environments/test/VARIABLE1 "envd_var1"
etcdctl set /environments/test/VARIABLE2 "envd_var2"

echo "Running tests..."
godep go test -v