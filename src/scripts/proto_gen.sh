#!/bin/sh

###### proto ####
awk -f proto.awk proto.txt > proto.go 
awk -f proto_func.awk proto.txt >> proto.go 

##### api generator #####
printf "package protos\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "import . \"types\"\n" >> api.go
printf "\n" >> api.go

awk -f api.awk api.txt >> api.go 
awk -f api_rcode.awk api.txt >> api.go 

printf "var ProtoHandler map[uint16]func(*Session, *packet.Packet) ([]byte, error) = map[uint16]func(*Session, *packet.Packet)([]byte, error){\n" >> api.go
awk -f api_bind.awk api.txt >> api.go 
printf "}" >> api.go

#### hub proto
awk -f proto.awk hub_proto.txt > hub_proto.go 
awk -f proto_func.awk hub_proto.txt >> hub_proto.go 

######## hub api ############
printf "package protos\n" > hub_api.go
printf "\n" >> hub_api.go
printf "import \"misc/packet\"\n" >> hub_api.go
printf "\n" >> hub_api.go

awk -f api.awk hub_api.txt >> hub_api.go 
awk -f api_rcode.awk hub_api.txt >> hub_api.go 

printf "var ProtoHandler map[uint16]func(int32, *packet.Packet) ([]byte, error) = map[uint16]func(int32, *packet.Packet)([]byte, error){\n" >> hub_api.go
awk -f api_bind.awk hub_api.txt >> hub_api.go 
printf "}" >> hub_api.go

#### move #################
mv proto.go ../agent/protos
mv api.go ../agent/protos

mv hub_proto.go ../hub/protos
mv hub_api.go ../hub/protos
