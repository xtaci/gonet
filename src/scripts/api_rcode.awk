###########################################################
## Scripts for generate protocol code(uint16)->(string)
##
BEGIN { RS = ""; FS ="\n" 
print "var RCode map[uint16]string = map[uint16]string {"
}
{
	name = ""
	for (i=1;i<=NF;i++)
	{
		if ($i ~ /^#.*/) {
			continue
		}
		split($i, a, ":")
		if (a[1] == "packet_type") {
			type = a[2]
		} else if (a[1] == "name") {
			name = a[2]
		} else if (a[1] == "payload") {
			payload = a[2]
		} else if (a[1] == "desc") {
			desc = a[2]
		}
	}
	if (name!= "") {
		print "\t"type":\""name"\","
	}
}
END {
print "}\n"	
}
