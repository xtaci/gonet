#!/bin/sh

##################################################
###   client proto & api
##################################################
printf "package protos\n" > proto.go
awk -f proto.awk proto.txt >> proto.go 
awk -f proto_func.awk proto.txt >> proto.go 

printf "package protos\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "import . \"types\"\n" >> api.go
printf "\n" >> api.go

awk -f api.awk api.txt >> api.go 
awk -f api_rcode.awk api.txt >> api.go 

printf "var ProtoHandler map[uint16]func(*Session, *packet.Packet) ([]byte, error) = map[uint16]func(*Session, *packet.Packet)([]byte, error){\n" >> api.go
awk -f api_bind_req.awk api.txt >> api.go 
printf "}" >> api.go

mv -f proto.go ../agent/client_protos
mv -f api.go ../agent/client_protos

##################################################
### hub proto & api
##################################################
printf "package protos\n" > proto.go
awk -f proto.awk hub_proto.txt >> proto.go 
awk -f proto_func.awk hub_proto.txt >> proto.go 

printf "package protos\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "\n" >> api.go

awk -f api.awk hub_api.txt >> api.go 
awk -f api_rcode.awk hub_api.txt >> api.go 

printf "var ProtoHandler map[uint16]func(int32, *packet.Packet) []byte = map[uint16]func(int32, *packet.Packet) []byte {\n" >> api.go
awk -f api_bind_req.awk hub_api.txt >> api.go 
printf "}" >> api.go

mv -f proto.go ../hub/protos
mv -f api.go ../hub/protos

##################################################
### event proto & api
##################################################
printf "package protos\n" > proto.go 
awk -f proto.awk event_proto.txt >> proto.go 
awk -f proto_func.awk event_proto.txt >> proto.go 

printf "package protos\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "\n" >> api.go

awk -f api.awk event_api.txt >> api.go 
awk -f api_rcode.awk event_api.txt >> api.go 

printf "var ProtoHandler map[uint16]func(*packet.Packet) []byte = map[uint16]func(*packet.Packet) []byte {\n" >> api.go
awk -f api_bind_req.awk event_api.txt >> api.go 
printf "}" >> api.go

mv -f proto.go ../event/protos
mv -f api.go ../event/protos
