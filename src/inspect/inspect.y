/* nex rp.nex && goyacc rp.y && 6g rp.nn.go y.go && 6l rp.nn.6 */
%{
package inspect
import "fmt"
import "net"
%}

%union { 
	n int;
	s string;
	conn net.Conn
} 

%token INSPECT 
%token VARS
%%
input:    /* empty */
       | input line
;

line:     '\n'
       | exp '\n'      { fmt.Fprintln(_conn, ">"); }
;

exp:    INSPECT VARS { fmt.Fprintln(_conn, $2.s); }
;
%%
