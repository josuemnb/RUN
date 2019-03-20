package run

type Return struct {
	Value Node
	Type  Type
}

func (p *Module) Return() Node {
	if !p.insideFunction() {
		p.error("Return only inside function", 0)
	}
	if p.ActualFunction.Return.Kind == VOID {
		if !p.match(EOL) {
			p.error("Expecting end of line", 1)
		}
		return Node{Type: RETURN, Value: Return{Type: *p.getTypeByKind(VOID)}}
	}
	t := p.assignment()
	typ := p.typeOf(t)
	if typ.Kind == 0 {
		p.error("Unknown identifier", 1)
	}
	if typ.Kind != QUOTE && p.ActualFunction.Return.Kind != STRING && typ.Kind != p.ActualFunction.Return.Kind {
		p.error("Types mismatches on returning", 1)
	}
	p.ActualFunction.Returned++
	p.consume(EOL, "Expecting end of line")
	return Node{Type: RETURN, Value: Return{Value: t, Type: typ}}
}

func (t *Transpiler) Return(node Node) {
	t.file.WriteString("return ")
	ret := node.Value.(Return)
	if ret.Type.Kind != VOID {
		t.Transpile(ret.Value)
	}
	t.file.WriteString(";\n")
}
