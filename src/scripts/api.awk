###########################################################
## Scripts for generate protocol string->code(uint16)
##
## packet_type:0
## name:heart_beat_req
## payload:null
## desc:心跳包..
##
BEGIN { RS = ""; FS ="\n" 
print "var Code = map[string]int16 {"
}
{
	for (i=1;i<=NF;i++)
	{
		if ($i ~ /^#.*/) {
			continue
		}

		split($i, a, ":")
		if (a[1] == "packet_type") {
			array["packet_type"] = a[2]
		} else if (a[1] == "name") {
			array["name"] = a[2]
		} else if (a[1] == "payload") {
			array["payload"] = a[2]
		} else if (a[1] == "desc") {
			array["desc"] = a[2]
		}
	}

	if ("packet_type" in array && "name" in array) {
		print "\t\""array["name"]"\":"array["packet_type"]",\t// "array["desc"]
	}

	delete array
}
END {
print "}\n"	
}
