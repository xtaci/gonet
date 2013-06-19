/* nex rp.nex && goyacc rp.y && 6g rp.nn.go y.go && 6l rp.nn.6 */
%{
package inspect
import (
	"fmt"
	"helper"
)
%}

%union { 
	n int;
	s string;
	nodes string;
} 

%token INSPECT 
%token GC 
%token FIELD
%token LIST 
%token ID 
%token QUIT
%token HELP 
%%
prog:	
		prog line
		|
		;

line:
		'\r' '\n'
		{ 
			prompt(conn)
		}
		| expr '\r' '\n'
		;

expr:
		list|help|quit|inspect|gc
		;

list:
		LIST 
		{
			ListAll(conn)
			prompt(conn)
		}
		;

help:
		HELP 
		{
			fmt.Fprintln(conn, "\t(p)rint user_id: inspect a user struct")
			fmt.Fprintln(conn, "\t(p)rint user_id.Field1.Field2...: dotted fields")
			fmt.Fprintln(conn, "\t(l)ist: list all online users")
			fmt.Fprintln(conn, "\tgc: force a garbage collection")
			prompt(conn) 
		}
		;

quit:
		QUIT 
		{ 
			conn.Close() 
		}
		;
	
inspect:	
		INSPECT ID 
		{ 
			Inspect(int32($2.n), conn)
			prompt(conn)
		}
		|
		INSPECT ID FIELD 
		{
			InspectField(int32($2.n), $3.nodes, conn) 
			prompt(conn)
		}
		;
gc:
		GC
		{
			fmt.Fprintln(conn,"before:")
			helper.FprintGCSummary(conn)
			helper.GC()
			fmt.Fprintln(conn,"after:")
			helper.FprintGCSummary(conn)
		}
		;
%%
