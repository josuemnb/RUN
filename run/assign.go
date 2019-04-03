package run

type Assign struct {
	Kind  Type
	Name  string
	Right Node
}

func isConst(s string) bool {
	for _, c := range s {
		if (c < 'A' || c > 'Z') && c != '_' {
			return false
		}
	}
	return true
}

func (p *Module) Assign() Node {
	n := p.advance()
	if _, ok := p.isVar(n.Lexeme); ok {
		p.rollBack()
		return p.assignment()
	}
	p.advance()
	e := p.assignment()
	k := p.typeOf(e)
	if k.Kind == QUOTE {
		k = *p.getTypeByKind(STRING)
	}
	if k.Kind == 0 {
		p.error("Unknown identifier", 1)
	}
	p.Scopes[p.CurScope][n.Lexeme] = Variable{Name: n.Lexeme, Type: k}
	return Node{Type: ASSIGN, Value: Assign{Kind: k, Name: n.Lexeme, Right: e}}
}

func (t *Transpiler) Assign(node Node) {
	assign := node.Value.(Assign)
	if isConst(assign.Name) {
		t.file.WriteString("const ")
	}
	if assign.Kind.Kind >= STRING || assign.Kind.Kind == QUOTE {
		t.file.WriteString(assign.Kind.Real + " var_" + assign.Name + "(")
		t.Transpile(assign.Right)
		t.file.WriteString(")")
	} else {
		t.file.WriteString(assign.Kind.Real + " var_" + assign.Name + "=")
		t.Transpile(assign.Right)
	}
}
