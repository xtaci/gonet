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
%token ID 
%token HELP 
%%
input:    /* empty */
       | input exp
;

exp:    INSPECT ID { Inspect(int32($2.n), conn); prompt(conn); }
		| HELP {fmt.Fprintln(conn, "\tinspect user_id"); prompt(conn);}
;
%%
