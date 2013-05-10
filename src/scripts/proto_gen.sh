#!/bin/sh

awk -f proto.awk proto.txt > proto.go 
awk -f proto_func.awk proto.txt >> proto.go 
awk -f api.awk api.txt > api.go 
awk -f api_rcode.awk api.txt >> api.go 
awk -f api_bind.awk api.txt >> api.go 
mv proto.go ../agent/protos
mv api.go ../agent/protos

#### hub
awk -f proto.awk hub_proto.txt > hub_proto.go 
awk -f proto_func.awk hub_proto.txt >> hub_proto.go 
awk -f api.awk hub_api.txt > hub_api.go 
awk -f api_rcode.awk hub_api.txt >> hub_api.go 
awk -f api_bind.awk hub_api.txt >> hub_api.go 
mv hub_proto.go ../hub/protos
mv hub_api.go ../hub/protos
