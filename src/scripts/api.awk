BEGIN { RS = ""; FS ="\n" 
print "package protos"
print ""
print "import \"misc/packet\"\n"
}
{

	for (i=1;i<=NF;i++)
	{
		if ($i ~ /^#.*/ || $i ~ /^===/) {
			continue
		}

		split($i, a, " ")
		if (a[1] ~ /.*=/) {
			print "type", substr(a[1], 1, length(a[1])-1), "struct {"
			typeok = "true"
		} else if (a[2] == "string") {
			print "\t"a[1] " string"
		} else if (a[2] == "integer") {
			print "\t"a[1] " int32"
		} else if (a[2] == "boolean") {
			print "\t"a[1] " byte"
		} else if (a[2] == "float") {
			print "\t"a[1] " float32"
		} else if (a[2] == "array") {
			print "\t"a[1]" []*"a[3]
		}
	}

	if (typeok) print "}\n"
}
END { }
