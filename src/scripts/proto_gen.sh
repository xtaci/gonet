#!/bin/sh

awk -f proto.awk proto.txt > proto.go 
awk -f proto_func.awk proto.txt >> proto.go 
awk -f api.awk api.txt > api.go 
awk -f api_rcode.awk api.txt >> api.go 
awk -f api_bind.awk api.txt >> api.go 
mv proto.go ../agent/protos
mv api.go ../agent/protos
