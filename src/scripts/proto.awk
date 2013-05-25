###########################################################
## generate proto payload struct 
##
BEGIN { RS = ""; FS ="\n" 
print ""
print "import \"misc/packet\"\n"
}
{

	typeok = false
	for (i=1;i<=NF;i++)
	{
		if ($i ~ /^#.*/ || $i ~ /^===/) {
			continue
		}

		split($i, a, " ")
		if (a[1] ~ /[A-Za-z_]+=/) {
			name = substr(a[1],1, match(a[1],/=/)-1)
			print "type",name, "struct {"
			typeok = "true"
		} else if (a[2] == "string") {
			print "\tF_"a[1] " string"
		} else if (a[2] == "integer" || a[2]=="int32") {
			print "\tF_"a[1] " int32"
		} else if (a[2] == "uint32") {
			print "\tF_"a[1] " uint32"
		} else if (a[2] == "long" || a[2]=="int64") {
			print "\tF_"a[1] " int64"
		} else if (a[2] == "uint64") {
			print "\tF_"a[1] " uint64"
		} else if (a[2] == "boolean") {
			print "\tF_"a[1] " byte"
		} else if (a[2] == "float") {
			print "\tF_"a[1] " float32"
		} else if (a[2] == "array") {
			print "\tF_"a[1]" []"a[3]
		} else {
			print "\tF_"a[1] " "a[2]
		}
	}

	if (typeok) print "}\n"
}
END { }
