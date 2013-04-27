BEGIN { RS = ""; FS ="\n" 
print "package protos"
print ""
print "import \"misc/packet\"\n"
print "import . \"types\"\n"
print ""
print "var Code map[string]uint16 = map[string]uint16 {"
}
{
	for (i=1;i<=NF;i++)
	{
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
	if (name != "") {
		print "\t\""name"\":"type",\t// payload:"payload" "desc
	}
}
END {
print "}\n"	
}
