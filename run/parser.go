package run

import (
	"os"
)

type Parser struct {
	insideClass    bool
	insideFunction bool
	insideLoop     bool
	tokens         []Token
	current        int
	returnType     Type
	length         int
	returned       int
	ignore         bool
	class          Class
	function       Function
	stmts          []Node
}

func (p *Module) Parse() []Node {
	stmts := make([]Node, 0)
	for !p.isAtEnd() {
		s := p.parse()
		if s.Type == EMPTY {
			continue
		}
		stmts = append(stmts, s)
	}
	return stmts
}

func (p *Module) error(e string, pos int) {
	if p.Tokens[p.CurToken].Type == EOL {
		p.CurToken--
	}
	println("Error at file '"+p.Name+" ':", e, "'"+p.Tokens[p.CurToken-pos].Lexeme+"'", "at line", p.Tokens[p.CurToken].Line)
	p.Finish(false)
	os.Exit(-1)
}

func (p *Module) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Module) peek() Token {
	return p.Tokens[p.CurToken]
}

func (p *Module) parse() Node {
	if p.match(EOL) {
		return Node{Type: EOL}
	}
	if p.match(CPP) {
		return p.Cpp()
	}
	if p.insideFunction() {
		if p.match(IF) {
			return p.ifStatment()
		}
		if p.match(PRINTLN) {
			return p.Print(true)
		}
		if p.match(PRINT) {
			return p.Print(false)
		}
		if p.match(RETURN) {
			return p.Return()
		}
		if p.match(LOOP) {
			return p.Loop()
		}
		if p.match(ASM) {
			return p.Asm()
		}
		if p.match(BREAK) {
			if p.InsideLoop == false {
				p.error("Unexpected break", 1)
			}
			p.consume(EOL, "Expecint end of line")
			return Node{Type: BREAK}
		}
		if p.test(DECLARE) {
			return p.Declare()
		}
		if p.test(EQUAL) {
			return p.Assign()
		}
		return p.assignment()
	}
	if p.insideClass() && p.insideFunction() == false {
		if p.match(EQUAL_EQUAL, EQUAL, BANG, BANG_EQUAL, LESS, LESS_EQUAL, GREATER, GREATER_EQUAL, MINUS_EQUAL, PLUS_EQUAL, MUL_EQUAL, DIV_EQUAL, PLUS, MINUS) {
			return p.operator()
		}
	}
	if p.insideClass() == false && p.insideFunction() == false {
		if p.match(MAIN) {
			return p.Main()
		}
		if p.test(LEFT_BRACE) {
			return p.Class(false)
		}
		if p.test(INTERFACE) {
			return p.Interface()
		}
		if p.test(EXTENDS) {
			return p.Class(true)
		}
	}
	if p.match(IMPORT) {
		return p.Import()
	}
	if p.match(MODULE) {
		return p.Module()
	}
	if p.test(DECLARE) {
		return p.Declare()
	}
	if p.test(LEFT_PAREN) {
		return p.Function()
	}
	p.error("Code outside function", 0)
	os.Exit(-3)
	return Node{}
}

func (p *Module) addVar(v string, t int) {
	if p.insideClass() {

	} else {
		p.Scopes[p.CurScope][v] = Variable{Name: v, Type: *p.getTypeByKind(t)}
	}
}

func (p *Module) isVar(n string) (Variable, bool) {
	if p.insideClass() {
		if v, ok := p.ActualClass.Fields[n]; ok {
			return v, true
		}
	}
	for i := p.CurScope; i >= 0; i-- {
		if v, ok := p.Scopes[i][n]; ok {
			return v, ok
		}

	}
	return Variable{}, false
}

func (p *Module) isFunction(f string) bool {
	if p.insideClass() {
		if _, ok := p.ActualClass.Methods[f]; ok {
			return true
		}
	}
	_, ok := p.Functions[f]
	return ok
}

func (m *Module) isArray(n string) bool {
	// for i := curScope; i >= 0; i-- {
	// 	if v, ok := scopes[i][n]; ok {
	// 		return v.Array != 0
	// 	}
	// }
	return false
}

func (p *Module) rollBack() {
	p.CurToken--
}

func (p *Module) compare(l, r Node) {
	if p.Ignore {
		return
	}
	t0 := p.typeOf(l)
	t1 := p.typeOf(r)
	if t0.Kind == NULL && (t1.Kind >= STRING || t1.Kind == QUOTE) {
		return
	}
	if t1.Kind == NULL && (t0.Kind >= STRING || t0.Kind == QUOTE) {
		return
	}
	if t0.Kind == QUOTE && t1.Kind == STRING {
		return
	} else if t1.Kind == QUOTE && t0.Kind == QUOTE {
		return
	}
	if t0.Kind != t1.Kind {
		p.error("Mismatches kinds", 1)
	}
}

func (p *Module) assignment() Node {
	var e Node
	e = p.or()
	if p.match(EQUAL) {
		return Node{Type: BINARY, Value: Binary{Left: e, Op: "=", Right: p.assignment()}}
	} else if p.match(PLUS_EQUAL) {
		t := p.typeOf(e)
		if t.Kind != STRING && t.Kind != NUMBER && t.Kind != REAL {
			p.error("Type not allowed for appending", 1)
		}
		return Node{Type: BINARY, Value: Binary{Left: e, Op: p.previous().Lexeme, Right: p.assignment()}}
	} else if p.match(MINUS_EQUAL) {
		t := p.typeOf(e)
		if t.Kind != NUMBER && t.Kind != REAL {
			p.error("Type not allowed for appending", 1)
		}
		return Node{Type: BINARY, Value: Binary{Left: e, Op: p.previous().Lexeme, Right: p.assignment()}}
	}
	return e
}

func (p *Module) match(kinds ...int) bool {
	for _, s := range kinds {
		if p.check(s) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Module) advance() Token {
	if !p.isAtEnd() {
		p.CurToken++
	}
	return p.previous()
}

func (p *Module) previous() Token {
	return p.Tokens[p.CurToken-1]
}

func (p *Module) checkAll(t ...int) bool {
	for _, s := range t {
		if p.check(s) {
			return true
		}
	}
	return false
}

func (p *Module) check(t int) bool {
	if p.isAtEnd() {
		p.error("Unexpected end of file", 0)
		return false
	}
	return p.peek().Type == t
}

func (p *Module) testAll(t ...int) bool {
	for _, s := range t {
		if p.test(s) {
			return true
		}
	}
	return false
}

func (p *Module) test(t int) bool {
	if p.CurToken+1 >= len(p.Tokens) {
		return false
	}
	if p.Tokens[p.CurToken+1].Type == EOF {
		return false
	}
	return p.Tokens[p.CurToken+1].Type == t
}

func (p *Module) consume(t int, s string) Token {
	if p.check(t) {
		return p.advance()
	}
	p.error("Error: "+s, 0)
	return Token{}
}
