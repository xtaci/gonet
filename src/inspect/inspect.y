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
} 

%token INSPECT 
%token GC 
%token LIST 
%token NUM
%token VARIABLE
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
		INSPECT
		{
			fmt.Fprintln(conn,"\t(p)rint what?")
			prompt(conn)
		}
		|
		INSPECT NUM
		{ 
			Inspect(int32($2.n), conn)
			prompt(conn)
		}
		|
		INSPECT NUM fields
		{
			InspectField(int32($2.n), $3.s, conn) 
			prompt(conn)
		}
		;

fields: 
		'.' VARIABLE fields
		{
			$$.s = "." + $2.s + $3.s
		}
		|
		'.' VARIABLE 
		{
			$$.s = "." + $2.s
		}
		;

gc:
		GC
		{
			helper.GC()
			helper.FprintGCSummary(conn)
			prompt(conn)
		}
		;
%%
