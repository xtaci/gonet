###########################################################
## generate proto payload struct 
##
BEGIN { RS = ""; FS ="\n" 
print "package protos"
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
		} else if (a[2] == "integer") {
			print "\tF_"a[1] " int32"
		} else if (a[2] == "long") {
			print "\tF_"a[1] " int64"
		} else if (a[2] == "boolean") {
			print "\tF_"a[1] " byte"
		} else if (a[2] == "float") {
			print "\tF_"a[1] " float32"
		} else if (a[2] == "array") {
			print "\tF_"a[1]" []"a[3]
		}
	}

	if (typeok) print "}\n"
}
END { }
