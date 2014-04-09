###########################################################
## Scripts for generate ProtoHandler map binding code
##
BEGIN { RS = ""; FS ="\n" }
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
			if (a[2] !~ /.*_req$/) {
				break
			} else {
				array["name"] = a[2]
			}
		}
	}

	if ("packet_type" in array && "name" in array) {
		print "\t"array["packet_type"]":P_"array["name"]","
	}

	delete array
}
