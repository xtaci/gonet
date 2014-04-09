#!/bin/sh

##################################################
###   client proto & api
##################################################
printf "package net\n" > proto.go
gawk -f proto.awk proto.txt >> proto.go 
gawk -f proto_func.awk proto.txt >> proto.go 

printf "package net\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "import . \"types\"\n" >> api.go
printf "\n" >> api.go

gawk -f api.awk api.txt >> api.go 
gawk -f api_rcode.awk api.txt >> api.go 

printf "var ProtoHandler map[int16]func(*Session, *packet.Packet) []byte\n" >> api.go
printf "func init() {\n" >> api.go
printf "ProtoHandler = map[int16]func(*Session, *packet.Packet) []byte {\n" >> api.go
gawk -f api_bind_req.awk api.txt >> api.go 
printf "}" >> api.go
printf "}" >> api.go

mv -f proto.go ../agent/net
mv -f api.go ../agent/net
go fmt ../agent/net
##################################################
###   a copy of proto into ipc
##################################################
printf "package ipc_service\n" > proto.go
gawk -f proto.awk proto.txt >> proto.go 
gawk -f proto_func.awk proto.txt >> proto.go 

mv -f proto.go ../agent/ipc_service
go fmt ../agent/ipc_service

##################################################
### hub proto & api
##################################################
printf "package protos\n" > proto.go
gawk -f proto.awk hub_proto.txt >> proto.go 
gawk -f proto_func.awk hub_proto.txt >> proto.go 

printf "package protos\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "\n" >> api.go

gawk -f api.awk hub_api.txt >> api.go 
gawk -f api_rcode.awk hub_api.txt >> api.go 

printf "var ProtoHandler map[int16]func(int32, *packet.Packet) []byte\n" >> api.go
printf "func init() {\n" >> api.go
printf "ProtoHandler = map[int16]func(int32, *packet.Packet) []byte {\n" >> api.go
gawk -f api_bind_req.awk hub_api.txt >> api.go 
printf "}" >> api.go
printf "}" >> api.go

mv -f proto.go ../hub/protos
mv -f api.go ../hub/protos
go fmt ../hub/protos

##################################################
### a copy of hub proto & api into agent/hub_client
##################################################
printf "package hub_client\n" > proto.go
gawk -f proto.awk hub_proto.txt >> proto.go
gawk -f proto_func.awk hub_proto.txt >> proto.go

printf "package hub_client\n" > api.go
printf "\n" >> api.go

gawk -f api.awk hub_api.txt >> api.go
gawk -f api_rcode.awk hub_api.txt >> api.go

mv -f proto.go ../agent/hub_client
mv -f api.go ../agent/hub_client
go fmt ../agent/hub_client

##################################################
### stats proto & api
##################################################
printf "package protos\n" > proto.go 
gawk -f proto.awk stats_proto.txt >> proto.go 
gawk -f proto_func.awk stats_proto.txt >> proto.go 

printf "package protos\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "\n" >> api.go

gawk -f api.awk stats_api.txt >> api.go 
gawk -f api_rcode.awk stats_api.txt >> api.go 

printf "var ProtoHandler map[int16]func(*packet.Packet) []byte\n" >> api.go
printf "func init() {\n" >> api.go
printf "ProtoHandler = map[int16]func(*packet.Packet) []byte {\n" >> api.go
gawk -f api_bind_req.awk stats_api.txt >> api.go 
printf "}" >> api.go
printf "}" >> api.go

mv -f proto.go ../stats/protos
mv -f api.go ../stats/protos
go fmt ../stats/protos

##################################################
###   a copy of proto into stats_client
##################################################
printf "package stats_client\n" > proto.go
gawk -f proto.awk stats_proto.txt >> proto.go
gawk -f proto_func.awk stats_proto.txt >> proto.go

printf "package stats_client\n" > api.go
printf "\n" >> api.go

gawk -f api.awk stats_api.txt >> api.go
gawk -f api_rcode.awk stats_api.txt >> api.go

mv -f proto.go ../agent/stats_client
mv -f api.go ../agent/stats_client
go fmt ../agent/stats_client
