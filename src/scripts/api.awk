###########################################################
## Scripts for generate protocol string->code(uint16)
##
BEGIN { RS = ""; FS ="\n" 
print "var Code map[string]int16 = map[string]int16 {"
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
		print "\t\""name"\":"type",\t// payload:"payload" "desc
	}
}
END {
print "}\n"	
}
