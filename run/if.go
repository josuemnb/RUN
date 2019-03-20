package run

type If struct {
	Cond Node
	Then []Node
	Else []Node
}

func isCondition(c string) bool {
	switch c {
	case "<", "<=", ">", ">=", "==", "!=":
		return true
	}
	return false
}

func (p *Module) ifStatment() Node {
	var i If
	t := IF
	i.Cond = p.assignment()
	b := i.Cond.Value.(Binary)
	if !isCondition(b.Op) {
		p.error("Conditional term invalid '"+b.Op+"'", 2)
	}
	p.consume(LEFT_BRACE, "Expecting begin of IF block")
	p.consume(EOL, "Expecting end of line")
	i.Then = p.block()
	if p.match(ELSE) {
		if p.check(IF) {
			t = ELSEIF
		} else {
			p.consume(LEFT_BRACE, "Expecting { for if statment")
			p.consume(EOL, "Expecting end of line")
			i.Else = p.block()
		}
	}
	return Node{Type: t, Value: i}
}

func (t *Transpiler) If(node Node) {
	i := node.Value.(If)
	t.file.WriteString("if(")
	t.Transpile(i.Cond)
	t.file.WriteString(") {\n")
	if i.Then != nil {
		for _, n := range i.Then {
			t.Transpile(n)
		}
	}
	if node.Type == ELSEIF {
		t.file.WriteString("} else ")
	} else {
		if i.Else != nil {
			t.file.WriteString("} else {\n")
			for _, e := range i.Else {
				t.Transpile(e)
			}
		}
		t.file.WriteString("}\n")
	}
}
