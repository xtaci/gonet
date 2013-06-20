###########################################################
## generate proto payload struct 
##
BEGIN { RS = "==="; FS ="\n" 
print ""
print "import \"misc/packet\"\n"
TYPES["byte"]="byte"
TYPES["short"]="int16"
TYPES["int16"]="int16"
TYPES["uint16"]="uint16"
TYPES["string"]="string"
TYPES["integer"]="int32"
TYPES["int32"]="int32"
TYPES["uint32"]="uint32"
TYPES["long"]="int64"
TYPES["int64"]="int64"
TYPES["uint64"]="uint64"
TYPES["bool"]="bool"
TYPES["boolean"]="bool"
TYPES["float"]="float32"
TYPES["float32"]="float32"
TYPES["float64"]="float64"
}
{
	for (i=1;i<=NF;i++)
	{
		if ($i ~ /^#.*/) {
			continue
		}

		if ($i ~ /[A-Za-z_]+=/) {
			name = substr($i,1, match($i,/=/)-1)
			print "type " name " struct {"
			typeok = "true"
		} else {
			v = _field($i)
			if (v) {
				print v
			}
		}

	}

	if (typeok) print "}\n"
	typeok=false

}
END { }

function _field(line) {
	split(line, a, " ")

	if (a[2] in TYPES) {
		return "\tF_"a[1] " " TYPES[a[2]]
	} else if (a[2] == "array") {
		if (a[3] in TYPES) {
			return "\tF_"a[1]" []" TYPES[a[3]]
		} else {
			return "\tF_"a[1]" []" a[3]
		}
	}
}
