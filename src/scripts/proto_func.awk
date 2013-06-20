###########################################################
## generate protocol packet reader
##
BEGIN { RS = ""; FS ="\n" 
TYPES["byte"]="ReadByte"
TYPES["short"]="ReadS16"
TYPES["int16"]="ReadS16"
TYPES["uint16"]="ReadU16"
TYPES["string"]="ReadString"
TYPES["integer"]="ReadS32"
TYPES["int32"]="ReadS32"
TYPES["uint32"]="ReadU32"
TYPES["long"]="ReadS64"
TYPES["int64"]="ReadS64"
TYPES["uint64"]="ReadU64"
TYPES["bool"]="ReadBool"
TYPES["boolean"]="ReadBool"
TYPES["float"]="ReadFloat32"
TYPES["float32"]="ReadFloat32"
TYPES["float64"]="ReadFloat64"
}
{
	for (i=1;i<=NF;i++)
	{
		if ($i ~ /^#.*/ || $i ~ /^===/) {
			continue
		}

		split($i, a, " ")
		if (a[1] ~ /[A-Za-z_]+=/) {
			name = substr(a[1],1, match(a[1],/=/)-1)
			print "func PKT_"name"(reader *packet.Packet)(tbl "name", err error){"
			typeok = "true"
		} else if (a[2] ==  "array") {
			if (a[3] == "byte") { 		## bytes
				print "\ttbl.F_"a[1]", err = reader.ReadBytes()"
				print "\tcheckErr(err)\n"
			} else {	## struct
				print "\tnarr := uint16(0)\n"
				print "\tnarr,err = reader.ReadU16()"
				print "\tcheckErr(err)\n"
				print "\ttbl.F_"a[1]"=make([]"a[3]",narr)"
				print "\tfor i:=0;i<int(narr);i++ {"
				print "\t\ttbl.F_"a[1]"[i], err = PKT_"a[3]"(reader)"
				print "\t\tcheckErr(err)\n"
				print "\t}\n"
			}
		}
		else {
			print "\ttbl.F_"a[1]",err = reader." TYPES[a[2]] "()"
			print "\tcheckErr(err)\n"
		}
	}

	if (typeok) {
		print "\treturn"
		print "}\n"
	}

	typeok=false
}
END { }
