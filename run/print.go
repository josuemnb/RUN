package run

type Print struct {
	Body    []Node
	Code    []Type
	newline bool
}

func (p *Module) Print(nl bool) Node {
	p.Ignore = true
	body := make([]Node, 0)
	code := make([]Type, 0)
	if p.match(LEFT_PAREN) {
		for {
			a := p.assignment()
			code = append(code, p.typeOf(a))
			body = append(body, a)
			if p.match(RIGHT_PAREN) {
				break
			}
			p.consume(COMMA, "Expecting ) or ,")
		}
	} else {
		p.consume(LEFT_BRACE, "Expecting ( or {")
		p.consume(EOL, "Expecting end of line")
		for !p.match(RIGHT_BRACE) {
			a := p.assignment()
			code = append(code, p.typeOf(a))
			body = append(body, a)
			if p.match(EOL) {
				if p.match(RIGHT_BRACE) {
					break
				}
				body = append(body, Node{Type: NEWLINE})
			} else {
				p.consume(COMMA, "Expecting } or , or newline")
			}
		}
	}
	p.Ignore = false
	p.consume(EOL, "Expecting end of line")
	return Node{Type: PRINT, Value: Print{Body: body, Code: code, newline: nl}}
}

func (t *Transpiler) Print(node Node) {
	t.Printing = true
	p := node.Value.(Print)
	t.file.WriteString("printf(\"")
	for _, n := range p.Code {
		t.file.WriteString(typeToRepr(n.Kind))
	}
	if p.newline {
		t.file.WriteString("\\n")
	}
	t.file.WriteString("\"")
	for i, n := range p.Body {
		t.file.WriteString(",")
		if p.Code[i].Kind == BOOL {
			t.file.WriteString("BOOL(")
		}
		t.Transpile(n)
		if p.Code[i].Kind == STRING {
			t.file.WriteString(".value")
		}
		if p.Code[i].Kind == BOOL {
			t.file.WriteString(")")
		}
	}
	t.file.WriteString(");\n")
	t.Printing = false
}
