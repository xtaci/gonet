BEGIN { RS = ""; FS ="\n" }
{
	for (i=1;i<=NF;i++)
	{
		if ($i ~ /^#.*/ || $i ~ /^===/) {
			continue
		}

		split($i, a, " ")
		if (a[1] ~ /[A-Za-z_]+=/) {
			name = substr(a[1],1, match(a[1],/=/)-1)
			print "func pktread_"name"(reader *packet.Packet)(tbl *"name", err error){"
			print "\ttbl = &"name"{}"
			typeok = "true"
		} else if (a[2] == "string") {
			print "\ttbl."a[1]",err = reader.ReadString()"
		} else if (a[2] == "integer") {
			print "\ttbl."a[1]",err = reader.ReadS32()"
			print "\tcheckErr(err)"
		} else if (a[2] == "boolean") {
			print "\ttbl."a[1]",err = reader.ReadByte()"
			print "\tcheckErr(err)"
		} else if (a[2] == "float") {
			print "\ttbl."a[1]",err = reader.ReadFloat()"
			print "\tcheckErr(err)"
		} else if (a[2] == "array") {
			print "\tnarr,err2 := reader.ReadU16()"
			print "\tcheckErr(err2)"
			print "\ttbl."a[1]"=make([]*"a[3]",narr)"
			print "\tfor i:=0;i<int(narr);i++ {"
			print "\t\ttbl."a[1]"[i], err = pktread_"a[3]"(reader)"
			print "\t}"
		}
	}

	if (typeok) {
		print "\treturn"
		print "}\n"
	}
}
END { }
