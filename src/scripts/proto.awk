BEGIN { RS = ""; FS ="\n" 
print "package protos"
print ""
print "var Code map[string]uint16 "
print "func init() { "
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
		print "\tCode[\""name"\"]="type"\t// payload:"payload" "desc
	}
}
END {
print "}"	
}
