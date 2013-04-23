#!/bin/sh

awk -f proto.awk proto.txt > proto.go 
awk -f api.awk api.txt > api.go 
awk -f api_func.awk api.txt >> api.go 
mv proto.go ../agent/protos
mv api.go ../agent/protos
