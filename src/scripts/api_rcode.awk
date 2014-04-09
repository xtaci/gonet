###########################################################
## Scripts for generate protocol code(uint16)->(string)
##
BEGIN { RS = ""; FS ="\n" 
print "var RCode = map[int16]string {"
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
		print "\t"array["packet_type"]":\""array["name"]"\",\t// "array["desc"]
	}

	delete array
}
END {
print "}\n"	
}
