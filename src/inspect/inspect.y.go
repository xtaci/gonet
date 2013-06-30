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
	yys int
	n   int
	s   string
}

const INSPECT = 57346
const GC = 57347
const LIST = 57348
const NUM = 57349
const VARIABLE = 57350
const QUIT = 57351
const HELP = 57352

var yyToknames = []string{
	"INSPECT",
	"GC",
	"LIST",
	"NUM",
	"VARIABLE",
	"QUIT",
	"HELP",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line ./src/inspect/inspect.y:105

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 19
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 23

var yyAct = []int{

	19, 13, 14, 10, 20, 18, 12, 11, 3, 21,
	15, 16, 17, 9, 8, 7, 6, 5, 4, 2,
	1, 0, 22,
}
var yyPact = []int{

	-1000, -3, -1000, -2, 0, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 5, -1000, -1000, -7, -9, -1000, -1000,
	1, -9, -1000,
}
var yyPgo = []int{

	0, 20, 19, 18, 17, 16, 15, 14, 13, 0,
}
var yyR1 = []int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	4, 5, 6, 7, 7, 7, 9, 9, 8,
}
var yyR2 = []int{

	0, 2, 0, 2, 3, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 2, 3, 3, 2, 1,
}
var yyChk = []int{

	-1000, -1, -2, 11, -3, -4, -5, -6, -7, -8,
	6, 10, 9, 4, 5, 12, 11, 7, 12, -9,
	13, 8, -9,
}
var yyDef = []int{

	2, -2, 1, 0, 0, 5, 6, 7, 8, 9,
	10, 11, 12, 13, 18, 3, 0, 14, 4, 15,
	0, 17, 16,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	12, 3, 3, 11, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 13,
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
		//line ./src/inspect/inspect.y:29
		{
			prompt(conn)
		}
	case 10:
		//line ./src/inspect/inspect.y:41
		{
			ListAll(conn)
			prompt(conn)
		}
	case 11:
		//line ./src/inspect/inspect.y:49
		{
			fmt.Fprintln(conn, "\t(p)rint user_id: inspect a user struct")
			fmt.Fprintln(conn, "\t(p)rint user_id.Field1.Field2...: dotted fields")
			fmt.Fprintln(conn, "\t(l)ist: list all online users")
			fmt.Fprintln(conn, "\tgc: force a garbage collection")
			prompt(conn)
		}
	case 12:
		//line ./src/inspect/inspect.y:60
		{
			conn.Close()
		}
	case 13:
		//line ./src/inspect/inspect.y:67
		{
			fmt.Fprintln(conn, "\t(p)rint what?")
			prompt(conn)
		}
	case 14:
		//line ./src/inspect/inspect.y:73
		{
			Inspect(int32(yyS[yypt-0].n), conn)
			prompt(conn)
		}
	case 15:
		//line ./src/inspect/inspect.y:79
		{
			InspectField(int32(yyS[yypt-1].n), yyS[yypt-0].s, conn)
			prompt(conn)
		}
	case 16:
		//line ./src/inspect/inspect.y:87
		{
			yyVAL.s = "." + yyS[yypt-1].s + yyS[yypt-0].s
		}
	case 17:
		//line ./src/inspect/inspect.y:92
		{
			yyVAL.s = "." + yyS[yypt-0].s
		}
	case 18:
		//line ./src/inspect/inspect.y:99
		{
			helper.GC()
			helper.FprintGCSummary(conn)
			prompt(conn)
		}
	}
	goto yystack /* stack new state and value */
}
