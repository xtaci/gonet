#!/bin/sh

##################################################
###   client proto & api
##################################################
printf "package ipc\n" > proto.go
gawk -f proto.awk proto.txt >> proto.go 
gawk -f proto_func.awk proto.txt >> proto.go 

printf "package ipc\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "import . \"types\"\n" >> api.go
printf "\n" >> api.go

gawk -f api.awk api.txt >> api.go 
gawk -f api_rcode.awk api.txt >> api.go 

printf "var ProtoHandler map[uint16]func(*Session, *packet.Packet) []byte = map[uint16]func(*Session, *packet.Packet) []byte {\n" >> api.go
gawk -f api_bind_req.awk api.txt >> api.go 
printf "}" >> api.go

mv -f proto.go ../agent/ipc
mv -f api.go ../agent/ipc

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

printf "var ProtoHandler map[uint16]func(int32, *packet.Packet) []byte = map[uint16]func(int32, *packet.Packet) []byte {\n" >> api.go
gawk -f api_bind_req.awk hub_api.txt >> api.go 
printf "}" >> api.go

mv -f proto.go ../hub/protos
mv -f api.go ../hub/protos

##################################################
### event proto & api
##################################################
printf "package protos\n" > proto.go 
gawk -f proto.awk event_proto.txt >> proto.go 
gawk -f proto_func.awk event_proto.txt >> proto.go 

printf "package protos\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "\n" >> api.go

gawk -f api.awk event_api.txt >> api.go 
gawk -f api_rcode.awk event_api.txt >> api.go 

printf "var ProtoHandler map[uint16]func(*packet.Packet) []byte = map[uint16]func(*packet.Packet) []byte {\n" >> api.go
gawk -f api_bind_req.awk event_api.txt >> api.go 
printf "}" >> api.go

mv -f proto.go ../event/protos
mv -f api.go ../event/protos

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

printf "var ProtoHandler map[uint16]func(*packet.Packet) []byte = map[uint16]func(*packet.Packet) []byte {\n" >> api.go
gawk -f api_bind_req.awk stats_api.txt >> api.go 
printf "}" >> api.go

mv -f proto.go ../stats/protos
mv -f api.go ../stats/protos
