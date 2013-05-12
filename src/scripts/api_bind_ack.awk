###########################################################
## Scripts for generate ProtoHandler map binding code
##
BEGIN { RS = ""; FS ="\n" }
{
	name = type=""
	for (i=1;i<=NF;i++)
	{
		split($i, a, ":")
		if (a[1] == "packet_type") {
			type = a[2]
		} else if (a[1] == "name") {
			if (a[2] !~ /.*_ack/) {
				break
			} else {
				name = a[2]
			}
		}
	}
	if (name != "" && type !="") {
		print "\t"type":P_"name","
	}
}
