//line ./src/inspect/inspect.y:2
package inspect

import __yyfmt__ "fmt"

//line ./src/inspect/inspect.y:2
import (
	"fmt"
	"helper"
)

//line ./src/inspect/inspect.y:9
type yySymType struct {
	yys   int
	n     int
	s     string
	nodes string
}

const INSPECT = 57346
const GC = 57347
const FIELD = 57348
const LIST = 57349
const ID = 57350
const QUIT = 57351
const HELP = 57352

var yyToknames = []string{
	"INSPECT",
	"GC",
	"FIELD",
	"LIST",
	"ID",
	"QUIT",
	"HELP",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line ./src/inspect/inspect.y:87

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 16
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 19

var yyAct = []int{

	13, 14, 18, 10, 15, 12, 11, 3, 16, 17,
	19, 9, 8, 7, 6, 5, 4, 2, 1,
}
var yyPact = []int{

	-1000, -4, -1000, -8, -3, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 1, -1000, -1000, -10, 4, -1000, -1000,
}
var yyPgo = []int{

	0, 18, 17, 16, 15, 14, 13, 12, 11,
}
var yyR1 = []int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	4, 5, 6, 7, 7, 8,
}
var yyR2 = []int{

	0, 2, 0, 2, 3, 1, 1, 1, 1, 1,
	1, 1, 1, 2, 3, 1,
}
var yyChk = []int{

	-1000, -1, -2, 11, -3, -4, -5, -6, -7, -8,
	7, 10, 9, 4, 5, 12, 11, 8, 12, 6,
}
var yyDef = []int{

	2, -2, 1, 0, 0, 5, 6, 7, 8, 9,
	10, 11, 12, 0, 15, 3, 0, 13, 4, 14,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	12, 3, 3, 11,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %U %s\n", uint(char), yyTokname(c))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf("saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 3:
		//line ./src/inspect/inspect.y:30
		{
			prompt(conn)
		}
	case 10:
		//line ./src/inspect/inspect.y:42
		{
			ListAll(conn)
			prompt(conn)
		}
	case 11:
		//line ./src/inspect/inspect.y:50
		{
			fmt.Fprintln(conn, "\t(p)rint user_id: inspect a user struct")
			fmt.Fprintln(conn, "\t(p)rint user_id.Field1.Field2...: dotted fields")
			fmt.Fprintln(conn, "\t(l)ist: list all online users")
			fmt.Fprintln(conn, "\tgc: force a garbage collection")
			prompt(conn)
		}
	case 12:
		//line ./src/inspect/inspect.y:61
		{
			conn.Close()
		}
	case 13:
		//line ./src/inspect/inspect.y:68
		{
			Inspect(int32(yyS[yypt-0].n), conn)
			prompt(conn)
		}
	case 14:
		//line ./src/inspect/inspect.y:74
		{
			InspectField(int32(yyS[yypt-1].n), yyS[yypt-0].nodes, conn)
			prompt(conn)
		}
	case 15:
		//line ./src/inspect/inspect.y:81
		{
			helper.GC()
			helper.FprintGCSummary(conn)
			prompt(conn)
		}
	}
	goto yystack /* stack new state and value */
}
