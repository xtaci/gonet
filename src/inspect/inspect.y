/* nex rp.nex && goyacc rp.y && 6g rp.nn.go y.go && 6l rp.nn.6 */
%{
package main
import "fmt"
%}

%union { 
	n int;
	s string;
} 

%token INSPECT 
%token VARS
%%
input:    /* empty */
       | input line
;

line:     '\n'
       | exp '\n'      { fmt.Println(">"); }
;

exp:    INSPECT VARS { fmt.Println($2.s); }
;
%%
