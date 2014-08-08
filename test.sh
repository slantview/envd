#!/bin/bash

echo "Killing off etcd if it is running..."
killall etcd

echo "Starting etcd..."
nohup etcd &
sleep 3

echo "Adding test variables..."
etcdctl set VARIABLE1 "envd_var1"
etcdctl set VARIABLE2 "envd_var2"

echo "Running tests..."
godep go test -v