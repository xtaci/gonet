/* nex rp.nex && goyacc rp.y && 6g rp.nn.go y.go && 6l rp.nn.6 */
%{
package inspect
import "fmt"
%}

%union { 
	n int;
	s string;
} 

%token INSPECT 
%token LIST 
%token ID 
%token CRLF
%token QUIT
%token HELP 
%%
input:	input exp CRLF
		|
		;

exp:	INSPECT ID { Inspect(int32($2.n), conn); prompt(conn)}
		| LIST {ListAll(conn); prompt(conn)}
		| HELP {fmt.Fprintln(conn, "\tinspect user_id: inspect a user")
				fmt.Fprintln(conn, "\tlist: list all online users")
				prompt(conn) }
		| QUIT { conn.Close() }
		;
%%
