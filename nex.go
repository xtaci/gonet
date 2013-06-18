// Substantial copy-and-paste from src/pkg/regexp.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)
import (
	"go/parser"
	"go/printer"
	"go/token"
)

type rule struct {
	regex             []rune
	code              string
	id, family, index int
}
type Error string

func (e Error) Error() string { return string(e) }

var (
	ErrInternal            = Error("internal error")
	ErrUnmatchedLpar       = Error("unmatched '('")
	ErrUnmatchedRpar       = Error("unmatched ')'")
	ErrUnmatchedLbkt       = Error("unmatched '['")
	ErrUnmatchedRbkt       = Error("unmatched ']'")
	ErrBadRange            = Error("bad range in character class")
	ErrExtraneousBackslash = Error("extraneous backslash")
	ErrBareClosure         = Error("closure applies to nothing")
	ErrBadBackslash        = Error("illegal backslash escape")
)

func ispunct(c rune) bool {
	for _, r := range "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~" {
		if c == r {
			return true
		}
	}
	return false
}

var escapes = []rune("abfnrtv")
var escaped = []rune("\a\b\f\n\r\t\v")

func escape(c rune) rune {
	for i, b := range escapes {
		if b == c {
			return escaped[i]
		}
	}
	return -1
}

type edge struct {
	// TODO: Make kind an enum (and use iota)
	// 2 = rune class edge
	// 1 = rune edge
	// 0 = empty edge
	// -1 = wild edge (".")
	kind   int
	r      rune   // Rune; for rune edges.
	lim    []rune // Pairs of limits for character class edges.
	negate bool   // True if the character class is negated.
	dst    *node  // Destination node.
}
type node struct {
	e      []*edge // Slice of outedges.
	n      int     // Index number. Scoped to a family.
	accept bool    // True if this is an accepting state.
	set    []int   // The NFA nodes represented by a DFA node.
}

func inClass(r rune, lim []rune) bool {
	for i := 0; i < len(lim); i += 2 {
		if lim[i] <= r && r <= lim[i+1] {
			return true
		}
	}
	return false
}
func gen(out io.Writer, x *rule) {
	// End rule
	if -1 == x.index {
		fmt.Fprintf(out, "a[%d].endcase = %d\n", x.family, x.id)
		return
	}
	s := x.regex
	// Regex -> NFA
	// We cannot have our alphabet be all Unicode characters. Instead,
	// we compute an alphabet for each regex:
	//
	//   1. Singles: we add single runes used in the regex: any rune not in a
	//   range. These are held in `sing`.
	//   2. Ranges: entire ranges become elements of the alphabet. If ranges in
	//   the same expression overlap, we break them up into non-overlapping
	//   ranges. The generated code checks singles before ranges, so there's no
	//   need to break up a range if it contains a single. These are maintained
	//   in sorted order in `lim`.
	//   3. Wild: we add an element representing all other runes.
	//
	// e.g. the alphabet of /[0-9]*[Ee][2-5]*/ is sing: { E, e },
	// lim: { [0-1], [2-5], [6-9] } and the wild element.
	sing := make(map[rune]bool)
	lim := make([]rune, 0, 8)
	var insertLimits func(l, r rune)
	// Insert a new range [l-r] into `lim`, breaking it up if it overlaps, and
	// discarding it if it coincides with an existing range. We keep `lim`
	// sorted.
	insertLimits = func(l, r rune) {
		var i int
		for i = 0; i < len(lim); i += 2 {
			if l <= lim[i+1] {
				break
			}
		}
		if len(lim) == i || r < lim[i] {
			lim = append(lim, 0, 0)
			copy(lim[i+2:], lim[i:])
			lim[i] = l
			lim[i+1] = r
			return
		}
		if l < lim[i] {
			lim = append(lim, 0, 0)
			copy(lim[i+2:], lim[i:])
			lim[i+1] = lim[i] - 1
			lim[i] = l
			insertLimits(lim[i], r)
			return
		}
		if l > lim[i] {
			lim = append(lim, 0, 0)
			copy(lim[i+2:], lim[i:])
			lim[i+1] = l - 1
			lim[i+2] = l
			insertLimits(l, r)
			return
		}
		// l == lim[i]
		if r == lim[i+1] {
			return
		}
		if r < lim[i+1] {
			lim = append(lim, 0, 0)
			copy(lim[i+2:], lim[i:])
			lim[i] = l
			lim[i+1] = r
			lim[i+2] = r + 1
			return
		}
		insertLimits(lim[i+1]+1, r)
	}
	pos := 0
	n := 0
	newNode := func() *node {
		res := new(node)
		res.n = n
		res.e = make([]*edge, 0, 8)
		n++
		return res
	}
	newEdge := func(u, v *node) *edge {
		res := new(edge)
		res.dst = v
		u.e = append(u.e, res)
		return res
	}
	newWildEdge := func(u, v *node) *edge {
		res := newEdge(u, v)
		res.kind = 1
		return res
	}
	newRuneEdge := func(u, v *node, r rune) *edge {
		res := newEdge(u, v)
		res.kind = 0
		res.r = r
		sing[r] = true
		return res
	}
	newNilEdge := func(u, v *node) *edge {
		res := newEdge(u, v)
		res.kind = -1
		return res
	}
	newClassEdge := func(u, v *node) *edge {
		res := newEdge(u, v)
		res.kind = 2
		res.lim = make([]rune, 0, 2)
		return res
	}
	nlpar := 0
	maybeEscape := func() rune {
		c := s[pos]
		if '\\' == c {
			pos++
			if len(s) == pos {
				panic(ErrExtraneousBackslash)
			}
			c = s[pos]
			switch {
			case ispunct(c):
			case escape(c) >= 0:
				c = escape(s[pos])
			default:
				panic(ErrBadBackslash)
			}
		}
		return c
	}
	pcharclass := func() (start, end *node) {
		start = newNode()
		end = newNode()
		e := newClassEdge(start, end)
		if len(s) > pos && '^' == s[pos] {
			e.negate = true
			pos++
		}
		isLowerLimit := true
		var left rune
		for {
			if len(s) == pos || s[pos] == ']' {
				if !isLowerLimit {
					panic(ErrBadRange)
				}
				return
			}
			switch s[pos] {
			case '-':
				panic(ErrBadRange)
			default:
				c := maybeEscape()
				pos++
				if len(s) == pos {
					panic(ErrBadRange)
				}
				switch {
				case isLowerLimit: // Lower limit.
					if '-' == s[pos] {
						pos++
						left = c
						isLowerLimit = false
					} else {
						e.lim = append(e.lim, c, c)
						sing[c] = true
					}
				case left <= c: // Upper limit.
					e.lim = append(e.lim, left, c)
					if left == c {
						sing[c] = true
					} else {
						insertLimits(left, c)
					}
					isLowerLimit = true
				default:
					panic(ErrBadRange)
				}
			}
		}
		panic("unreachable")
	}
	var pre func() (start, end *node)
	pterm := func() (start, end *node) {
		if len(s) == pos || s[pos] == '|' {
			end = newNode()
			start = end
			return
		}
		switch s[pos] {
		case '*', '+', '?':
			panic(ErrBareClosure)
		case ')':
			if 0 == nlpar {
				panic(ErrUnmatchedRpar)
			}
			end = newNode()
			start = end
			return
		case '(':
			nlpar++
			pos++
			start, end = pre()
			if len(s) == pos || ')' != s[pos] {
				panic(ErrUnmatchedLpar)
			}
		case '.':
			start = newNode()
			end = newNode()
			newWildEdge(start, end)
		case ']':
			panic(ErrUnmatchedRbkt)
		case '[':
			pos++
			start, end = pcharclass()
			if len(s) == pos || ']' != s[pos] {
				panic(ErrUnmatchedLbkt)
			}
		default:
			start = newNode()
			end = newNode()
			newRuneEdge(start, end, maybeEscape())
		}
		pos++
		return
	}
	pclosure := func() (start, end *node) {
		start, end = pterm()
		if start == end {
			return
		}
		if len(s) == pos {
			return
		}
		switch s[pos] {
		case '*':
			newNilEdge(end, start)
			nend := newNode()
			newNilEdge(end, nend)
			start = end
			end = nend
		case '+':
			newNilEdge(end, start)
			nend := newNode()
			newNilEdge(end, nend)
			end = nend
		case '?':
			newNilEdge(start, end)
		default:
			return
		}
		pos++
		return
	}
	pcat := func() (start, end *node) {
		for {
			nstart, nend := pclosure()
			if start == nil {
				start, end = nstart, nend
			} else if nstart != nend {
				end.e = make([]*edge, len(nstart.e))
				copy(end.e, nstart.e)
				end = nend
			}
			if nstart == nend {
				return
			}
		}
		panic("unreachable")
	}
	pre = func() (start, end *node) {
		start, end = pcat()
		for {
			if len(s) == pos {
				return
			}
			if s[pos] != '|' {
				return
			}
			pos++
			nstart, nend := pcat()
			tmp := newNode()
			newNilEdge(tmp, start)
			newNilEdge(tmp, nstart)
			start = tmp
			tmp = newNode()
			newNilEdge(end, tmp)
			newNilEdge(nend, tmp)
			end = tmp
		}
		panic("unreachable")
	}
	start, end := pre()
	end.accept = true

	// Compute shortlist of nodes, as we may have discarded nodes left over
	// from parsing. Also, make short[0] the start node.
	short := make([]*node, 0, n)
	{
		var visit func(*node)
		mark := make([]bool, n)
		newn := make([]int, n)
		visit = func(u *node) {
			mark[u.n] = true
			newn[u.n] = len(short)
			short = append(short, u)
			for _, e := range u.e {
				if !mark[e.dst.n] {
					visit(e.dst)
				}
			}
		}
		visit(start)
		for _, v := range short {
			v.n = newn[v.n]
		}
	}
	n = len(short)

	/*
		  {
		  var show func(*node)
		  mark := make([]bool, n)
		  show = func(u *node) {
		    mark[u.n] = true
		    print(u.n, ": ")
		    for _,e := range u.e {
		      print("(", e.kind, " ", e.r, ")")
		      print(e.dst.n)
		    }
		    println()
		    for _,e := range u.e {
		      if !mark[e.dst.n] {
			show(e.dst)
		      }
		    }
		  }
		  show(start)
		  }
	*/

	// NFA -> DFA
	nilClose := func(st []bool) {
		mark := make([]bool, n)
		var do func(int)
		do = func(i int) {
			v := short[i]
			for _, e := range v.e {
				if -1 == e.kind && !mark[e.dst.n] {
					st[e.dst.n] = true
					do(e.dst.n)
				}
			}
		}
		for i := 0; i < n; i++ {
			if st[i] && !mark[i] {
				mark[i] = true
				do(i)
			}
		}
	}
	todo := make([]*node, 0, n)
	tab := make(map[string]*node)
	buf := make([]byte, 0, 8)
	dfacount := 0
	{
		for i := 0; i < n; i++ {
			buf = append(buf, '0')
		}
		tmp := new(node)
		tmp.n = -1
		tab[string(buf)] = tmp
	}
	newDFANode := func(st []bool) (res *node, found bool) {
		buf = buf[:0]
		accept := false
		for i, v := range st {
			if v {
				buf = append(buf, '1')
				accept = accept || short[i].accept
			} else {
				buf = append(buf, '0')
			}
		}
		res, found = tab[string(buf)]
		if !found {
			res = new(node)
			res.set = make([]int, 0, 8)
			res.n = dfacount
			res.accept = accept
			dfacount++
			for i, v := range st {
				if v {
					res.set = append(res.set, i)
				}
			}
			tab[string(buf)] = res
		}
		return res, found
	}

	get := func(states []bool) *node {
		nilClose(states)
		node, old := newDFANode(states)
		if !old {
			todo = append(todo, node)
		}
		return node
	}
	states := make([]bool, n)
	states[0] = true
	get(states)
	for len(todo) > 0 {
		v := todo[len(todo)-1]
		todo = todo[0 : len(todo)-1]
		// Singles.
		for r, _ := range sing {
			states := make([]bool, n)
			for _, i := range v.set {
				for _, e := range short[i].e {
					if e.kind == 0 && e.r == r {
						states[e.dst.n] = true
					} else if e.kind == 1 {
						states[e.dst.n] = true
					} else if e.kind == 2 && e.negate != inClass(r, e.lim) {
						states[e.dst.n] = true
					}
				}
			}
			newRuneEdge(v, get(states), r)
		}
		// Character ranges.
		for j := 0; j < len(lim); j += 2 {
			states := make([]bool, n)
			for _, i := range v.set {
				for _, e := range short[i].e {
					if e.kind == 1 {
						states[e.dst.n] = true
					} else if e.kind == 2 && e.negate != inClass(lim[j], e.lim) {
						states[e.dst.n] = true
					}
				}
			}
			e := newClassEdge(v, get(states))
			e.lim = append(e.lim, lim[j], lim[j+1])
		}
		// Other.
		states := make([]bool, n)
		for _, i := range v.set {
			for _, e := range short[i].e {
				if e.kind == 1 || (e.kind == 2 && e.negate) {
					states[e.dst.n] = true
				}
			}
		}
		newWildEdge(v, get(states))
	}
	n = dfacount
	// DFA -> Go
	// TODO: Literal arrays instead of a series of assignments.
	fmt.Fprintf(out, "{\nvar acc [%d]bool\nvar fun [%d]func(rune) int\n", n, n)
	for _, v := range tab {
		if -1 == v.n {
			continue
		}
		if v.accept {
			fmt.Fprintf(out, "acc[%d] = true\n", v.n)
		}
		fmt.Fprintf(out, "fun[%d] = func(r rune) int {\n", v.n)
		fmt.Fprintf(out, "  switch(r) {\n")
		for _, e := range v.e {
			m := e.dst.n
			if e.kind == 0 {
				fmt.Fprintf(out, "  case %d: return %d\n", e.r, m)
			}
		}
		fmt.Fprintf(out, "  default:\n    switch {\n")
		for _, e := range v.e {
			m := e.dst.n
			if e.kind == 2 {
				fmt.Fprintf(out, "    case %d <= r && r <= %d: return %d\n",
					e.lim[0], e.lim[1], m)
			} else if e.kind == 1 {
				fmt.Fprintf(out, "    default: return %d\n", m)
			}
		}
		fmt.Fprintf(out, "    }\n  }\n  panic(\"unreachable\")\n}\n")
	}
	fmt.Fprintf(out, "a%d[%d].acc = acc[:]\n", x.family, x.index)
	fmt.Fprintf(out, "a%d[%d].f = fun[:]\n", x.family, x.index)
	fmt.Fprintf(out, "a%d[%d].id = %d\n", x.family, x.index, x.id)
	fmt.Fprintf(out, "}\n")
}

var standalone *bool
var customError *bool

func writeLex(out *bufio.Writer, rules []*rule) {
	if !*customError {
		out.WriteString(`func (yylex Lexer) Error(e string) {
  panic(e)
}`)
	}
	out.WriteString(`
func (yylex Lexer) Lex(lval *yySymType) int {
  for !yylex.isDone() {
    switch yylex.nextAction() {
    case -1:`)
	for _, x := range rules {
		fmt.Fprintf(out, "\n    case %d:  //%s/\n", x.id, string(x.regex))
		out.WriteString(x.code)
	}
	out.WriteString("    }\n  }\n  return 0\n}\n")
}
func writeNNFun(out *bufio.Writer, rules []*rule) {
	out.WriteString(`func(yylex Lexer) {
  for !yylex.isDone() {
    switch yylex.nextAction() {
    case -1:`)
	for _, x := range rules {
		fmt.Fprintf(out, "\n    case %d:  //%s/\n", x.id, string(x.regex))
		out.WriteString(x.code)
	}
	out.WriteString("    }\n  }\n  }")
}
func process(output io.Writer, input io.Reader) {
	in := bufio.NewReader(input)
	out := bufio.NewWriter(output)
	var r rune
	done := false
	regex := make([]rune, 0, 8)
	read := func() {
		var er error
		r, _, er = in.ReadRune()
		if er == io.EOF {
			done = true
		} else if er != nil {
			panic(er.Error())
		}
	}
	skipws := func() {
		for {
			read()
			if done {
				break
			}
			if strings.IndexRune(" \n\t\r", r) == -1 {
				break
			}
		}
	}
	var rules []*rule
	usercode := false
	familyn := 1
	id := 0
	newRule := func(family, index int) *rule {
		x := new(rule)
		rules = append(rules, x)
		x.family = family
		x.id = id
		x.index = index
		id++
		return x
	}
	buf := make([]rune, 0, 8)
	readCode := func() string {
		if '{' != r {
			panic("expected {")
		}
		buf = buf[:0]
		nesting := 1
		for {
			buf = append(buf, r)
			read()
			if done {
				panic("unmatched {")
			}
			if '{' == r {
				nesting++
			}
			if '}' == r {
				nesting--
				if 0 == nesting {
					break
				}
			}
		}
		buf = append(buf, r)
		return string(buf)
	}
	var decls string
	var parse func(int)
	parse = func(family int) {
		rulen := 0
		declvar := func() {
			decls += fmt.Sprintf("var a%d [%d]dfa\n", family, rulen)
		}
		for !done {
			skipws()
			if done {
				break
			}
			regex = regex[:0]
			if '>' == r {
				if 0 == family {
					panic("unmatched >")
				}
				x := newRule(family, -1)
				x.code = "yylex = yylex.pop()\n"
				declvar()
				skipws()
				x.code += readCode()
				return
			}
			delim := r
			read()
			if done {
				panic("unterminated pattern")
			}
			for {
				if r == delim && (len(regex) == 0 || regex[len(regex)-1] != '\\') {
					break
				}
				if '\n' == r {
					panic("regex interrupted by newline")
				}
				regex = append(regex, r)
				read()
				if done {
					panic("unterminated pattern")
				}
			}
			if "" == string(regex) {
				usercode = true
				break
			}
			skipws()
			if done {
				panic("last pattern lacks action")
			}
			x := newRule(family, rulen)
			rulen++
			x.regex = make([]rune, len(regex))
			copy(x.regex, regex)
			nested := false
			if '<' == r {
				skipws()
				if done {
					panic("'<' lacks action")
				}
				x.code = fmt.Sprintf("yylex = yylex.push(%d)\n", familyn)
				nested = true
			}
			x.code += readCode()
			if nested {
				familyn++
				parse(familyn - 1)
			}
		}
		if 0 != family {
			panic("unmatched <")
		}
		x := newRule(family, -1)
		x.code = "// [END]\n"
		declvar()
	}
	parse(0)

	if !usercode {
		return
	}

	skipws()
	buf = buf[:0]
	for !done {
		buf = append(buf, r)
		read()
	}
	fs := token.NewFileSet()
	t, err := parser.ParseFile(fs, "", string(buf), parser.ImportsOnly)
	if err != nil {
		panic(err.Error())
	}
	printer.Fprint(out, fs, t)

	var file *token.File
	fs.Iterate(func(f *token.File) bool {
		file = f
		return true
	})

	for m := file.LineCount(); m > 1; m-- {
		i := 0
		for '\n' != buf[i] {
			i++
		}
		buf = buf[i+1:]
	}

	fmt.Fprintf(out, `import ("bufio";"io";"strings")
type dfa struct {
  acc []bool
  f []func(rune) int
  id int
}
type family struct {
  a []dfa
  endcase int
}
`)
	out.WriteString(decls)
	out.WriteString("var a []family\n")
	out.WriteString("func init() {\n")

	fmt.Fprintf(out, "a = make([]family, %d)\n", familyn)
	for _, x := range rules {
		gen(out, x)
	}
	for i := 0; i < familyn; i++ {
		fmt.Fprintf(out, "a[%d].a = a%d[:]\n", i, i)
	}

	out.WriteString(`}
func getAction(c *frame) int {
  if -1 == c.match { return -1 }
  c.action = c.fam.a[c.match].id
  c.match = -1
  return c.action
}
type frame struct {
  atEOF bool
  action, match, matchn, n int
  buf []rune
  text string
  in *bufio.Reader
  state []int
  fam family
}
func newFrame(in *bufio.Reader, index int) *frame {
  f := new(frame)
  f.buf = make([]rune, 0, 128)
  f.in = in
  f.match = -1
  f.fam = a[index]
  f.state = make([]int, len(f.fam.a))
  return f
}
type Lexer []*frame
func NewLexer(in io.Reader) Lexer {
  stack := make([]*frame, 0, 4)
  stack = append(stack, newFrame(bufio.NewReader(in), 0))
  return stack
}
func (stack Lexer) isDone() bool {
  return 1 == len(stack) && stack[0].atEOF
}
func (stack Lexer) nextAction() int {
  c := stack[len(stack) - 1]
  for {
    if c.atEOF { return c.fam.endcase }
    if c.n == len(c.buf) {
      r,_,er := c.in.ReadRune()
      switch er {
      case nil: c.buf = append(c.buf, r)
      case io.EOF:
	c.atEOF = true
	if c.n > 0 {
	  c.text = string(c.buf)
	  return getAction(c)
	}
	return c.fam.endcase
      default: panic(er.Error())
      }
    }
    jammed := true
    r := c.buf[c.n]
    for i, x := range c.fam.a {
      if -1 == c.state[i] { continue }
      c.state[i] = x.f[c.state[i]](r)
      if -1 == c.state[i] { continue }
      jammed = false
      if x.acc[c.state[i]] {
	if -1 == c.match || c.matchn < c.n+1 || c.match > i {
	  c.match = i
	  c.matchn = c.n+1
	}
      }
    }
    if jammed {
      a := getAction(c)
      if -1 == a { c.matchn = c.n + 1 }
      c.n = 0
      for i, _ := range c.state { c.state[i] = 0 }
      c.text = string(c.buf[:c.matchn])
      copy(c.buf, c.buf[c.matchn:])
      c.buf = c.buf[:len(c.buf) - c.matchn]
      return a
    }
    c.n++
  }
  panic("unreachable")
}
func (stack Lexer) push(index int) Lexer {
  c := stack[len(stack) - 1]
  return append(stack,
      newFrame(bufio.NewReader(strings.NewReader(c.text)), index))
}
func (stack Lexer) pop() Lexer {
  return stack[:len(stack) - 1]
}
func (stack Lexer) Text() string {
  c := stack[len(stack) - 1]
  return c.text
}
`)
	if !*standalone {
		writeLex(out, rules)
		out.WriteString(string(buf))
		out.Flush()
		return
	}
	m := 0
	const funmac = "NN_FUN"
	for m < len(buf) {
		m++
		if funmac[:m] != string(buf[:m]) {
			out.WriteString(string(buf[:m]))
			buf = buf[m:]
			m = 0
		} else if funmac == string(buf[:m]) {
			writeNNFun(out, rules)
			buf = buf[m:]
			m = 0
		}
	}
	out.WriteString(string(buf))
	out.Flush()
}

func main() {
	standalone = flag.Bool("s", false,
		`standalone code; no Lex() method, NN_FUN macro substitution`)
	customError = flag.Bool("e", false,
		`custom error func; no Error() method`)
	flag.Parse()
	if 0 == flag.NArg() {
		process(os.Stdout, os.Stdin)
		return
	}
	if flag.NArg() > 1 {
		println("nex: extraneous arguments after", flag.Arg(1))
		os.Exit(1)
	}
	if strings.HasSuffix(flag.Arg(1), ".go") {
		println("nex: input filename ends with .go:", flag.Arg(1))
		os.Exit(1)
	}
	basename := flag.Arg(0)
	n := strings.LastIndex(basename, ".")
	if n >= 0 {
		basename = basename[:n]
	}
	name := basename + ".nn.go"
	f, er := os.Open(flag.Arg(0))
	if f == nil {
		println("nex:", er.Error())
		os.Exit(1)
	}
	outf, er := os.Create(name)
	if outf == nil {
		println("nex:", er.Error())
		os.Exit(1)
	}
	process(outf, f)
	outf.Close()
	f.Close()
}
