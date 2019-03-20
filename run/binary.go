package run

type Binary struct {
	Left  Node
	Right Node
	Op    string
}

func (p *Module) or() Node {
	var e Node
	e = p.and()
	for p.match(OR) {
		right := p.equality()
		p.compare(e, right)
		e = Node{Type: BINARY, Value: Binary{e, right, "||"}}
	}
	return e
}

func (p *Module) and() Node {
	var e Node
	e = p.equality()
	for p.match(AND) {
		right := p.equality()
		p.compare(e, right)
		e = Node{Type: BINARY, Value: Binary{e, right, "&&"}}
	}
	return e
}

func (p *Module) equality() Node {
	var e Node
	e = p.comparation()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		op := p.previous()
		right := p.comparation()
		p.compare(e, right)
		e = Node{Type: BINARY, Value: Binary{e, right, op.Lexeme}}
	}
	return e
}

func (p *Module) comparation() Node {
	var e Node
	e = p.adition()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		op := p.previous()
		right := p.adition()
		p.compare(e, right)
		e = Node{Type: BINARY, Value: Binary{e, right, op.Lexeme}}
	}
	return e
}

func (p *Module) adition() Node {
	var e Node
	e = p.multiplication()
	for p.match(MINUS, PLUS) {
		op := p.previous()
		right := p.multiplication()
		p.compare(e, right)
		e = Node{Type: BINARY, Value: Binary{e, right, op.Lexeme}}
	}
	return e
}

func (p *Module) multiplication() Node {
	var e Node
	e = p.unary()
	for p.match(SLASH, STAR) {
		op := p.previous()
		right := p.unary()
		p.compare(e, right)
		e = Node{Type: BINARY, Value: Binary{e, right, op.Lexeme}}
	}
	return e
}

func (t *Transpiler) Binary(node Node) {
	b := node.Value.(Binary)
	t.Transpile(b.Left)
	t.file.WriteString(b.Op)
	t.Transpile(b.Right)
}
