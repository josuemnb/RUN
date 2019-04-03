package run

type Literal struct {
	Type  Type
	Value string
}

type Identifier struct {
	Name string
	Kind int
	Type Type
}

type Sequence struct {
	Left  Node
	Right Node
	Op    string
}

type Dot struct {
	Left Node
}

type Grouping struct {
	Group Node
}

type Bracket struct {
	Left  Identifier
	Index Node
}

type Between struct {
	Left  Node
	Right Node
}

func (p *Module) child(class Class, e Node) Node {
	var b Binary
	t := p.consume(IDENTIFIER, "Expecting identifier").Lexeme
	if p.match(LEFT_PAREN) {
		if class.isMethod(t) == false {
			p.error("Class '"+class.Name+"' doesnt contains method", 2)
		}
		real := p.CheckParams()
		f := class.Methods[t+real]
		if f.Name == "" {
			p.error("Method '"+t+"' not found", 2)
		}
		if f.Protection != PUBLIC {
			p.error("Class '"+class.Name+"' not accessible", 2)
		}
		b.Left = e
		b.Right = p.GetParams(f, METHOD)
		if p.match(DOT) {
			b.Right = Node{Type: DOT, Value: Dot{Left: b.Right}}
			return p.child(f.Return.Class, Node{Type: BINARY, Value: b})
		}
	} else if class.isField(t) {
		f := class.getField(t)
		if f.Protection != PUBLIC {
			p.error("Class '"+class.Name+"' not accessible", 1)
		}
	}
	return Node{Type: BINARY, Value: b}
}

func (p *Module) identifier() Node {
	return Node{Type: IDENTIFIER, Value: Identifier{Name: p.previous().Lexeme}}
}

func (p *Module) primary() Node {
	if p.match(FALSE, TRUE) {
		return Node{Type: LITERAL, Value: Literal{Type: *p.getTypeByKind(BOOL), Value: p.previous().Lexeme}}
	}
	if p.match(NIL) {
		return Node{Type: NULL}
	}
	if p.match(NUMBER, QUOTE, REAL) {
		return Node{Type: LITERAL, Value: Literal{Type: *p.getTypeByKind(p.previous().Type), Value: p.previous().Lexeme}}
	}
	if p.match(SUPER) {

	}
	if p.match(THIS) {
		return p.This()
	}
	if p.match(IDENTIFIER) {
		return p.identifier()
	}
	if p.match(LEFT_PAREN) {
		var e Node
		e = p.assignment()
		p.consume(RIGHT_PAREN, "Expecting )")
		return Node{Type: GROUPING, Value: Grouping{Group: e}}

	} else if p.match(EOL) {
		if p.Ignore {
			return Node{Type: NEWLINE}
		}
		return Node{Type: EOL}
	}
	p.error("Unrecognized Token", 0)
	return Node{Type: EMPTY}
}

func (t *Transpiler) Grouping(node Node) {
	t.file.WriteString("(")
	t.Transpile(node.Value.(Grouping).Group)
	t.file.WriteString(")")
}

func (t *Transpiler) Identifier(node Node) {
	i := node.Value.(Identifier)
	switch i.Kind {
	case FIELD:
		t.file.WriteString("field_" + i.Name)
	case VARIABLE:
		if i.Type.IsInterface {
			t.file.WriteString("inter_" + i.Name)
		} else {
			t.file.WriteString("var_" + i.Name)
		}
	case METHOD:
		t.file.WriteString("method_" + i.Name)
	case PARAM:
		t.file.WriteString("param_" + i.Name)
	case FUNCTION:
		t.file.WriteString("func_" + i.Name)
	case CLASS:
		t.file.WriteString("class_" + i.Name)
	default:
		// t.file.WriteString(i.Name)
	}
}

func (t *Transpiler) Literal(node Node) {
	l := node.Value.(Literal)
	if l.Type.Kind == QUOTE {
		t.file.WriteString("\"" + l.Value + "\"")
	} else {
		t.file.WriteString(l.Value)
	}
}

func (t *Transpiler) Dot(node Node) {
	d := node.Value.(Dot)
	t.Transpile(d.Left)
	t.file.WriteString(".")
}

func (t *Transpiler) Bracket(node Node) {
	b := node.Value.(Bracket)
	if b.Index.Type != BETWEEN {
		t.Transpile(Node{Type: IDENTIFIER, Value: b.Left})
		t.file.WriteString("[")
		t.Transpile(b.Index)
		t.file.WriteString("]")
		if b.Left.Type.Kind&STRING == STRING {
			t.file.WriteString(".value")
		}
	} else {
		t.file.WriteString("var_" + b.Left.Name + ".substring(")
		t.Transpile(b.Index)
		t.file.WriteString(")")
	}
}

func (t *Transpiler) Between(node Node) {
	b := node.Value.(Between)
	t.Transpile(b.Left)
	t.file.WriteString(",")
	t.Transpile(b.Right)
}
