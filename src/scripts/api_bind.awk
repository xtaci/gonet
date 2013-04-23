BEGIN { RS = ""; FS ="\n" 
print ""
print "var ProtoHandler map[uint16]func(*Session, *packet.Packet) ([]byte, error) = map[uint16]func(*Session, *packet.Packet)([]byte, error){"
}
{
	name = type=""
	for (i=1;i<=NF;i++)
	{
		split($i, a, ":")
		if (a[1] == "packet_type") {
			type = a[2]
		} else if (a[1] == "name") {
			if (a[2] !~ /.*_req/) {
				break
			} else {
				name = a[2]
			}
		}
	}
	if (name != "" && type !="") {
		print "\t"type":_"name","
	}
}
END {
print "}"	
}
