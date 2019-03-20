package run

type Function struct {
	Params     []Param
	Body       []Node
	Return     Type
	Name       string
	Protection Protection
	Returned   int
	Real       string
}

func (f *Function) isParam(n string) bool {
	for _, p := range f.Params {
		if p.Name == n {
			return true
		}
	}
	return false
}

func (f *Function) getParam(n string) Param {
	for _, p := range f.Params {
		if p.Name == n {
			return p
		}
	}
	return Param{}
}

func (p *Module) Function() Node {
	var f Function
	t := p.advance()
	if p.isFunction(t.Lexeme) {
		p.rollBack()
		return p.Call(p.Functions[t.Lexeme])
	}
	if _, ok := p.isVar(t.Lexeme); ok {
		p.error("Name already assigned to var", 0)
	}
	f.Name = t.Lexeme
	if p.Type == MODULE {
		f.Real = p.Name + "_"
	}
	f.Real += f.Name
	p.advance()
	p.BeginScope()
	f.Params = make([]Param, 0)
	if !p.match(RIGHT_PAREN) {
		for {
			var param Param
			n := p.advance()
			if f.isParam(n.Lexeme) {
				p.error("Name already assigned", 0)
			}
			param.Name = n.Lexeme
			if p.check(DECLARE) {
				cls, array := p.TypeDeclare(true)
				param.Type = cls
				p.Scopes[p.CurScope][param.Name] = Variable{Name: param.Name, Type: cls, Array: array}
			} else {
				param.Type = *p.getTypeByName("string")
				p.Scopes[p.CurScope][param.Name] = Variable{Name: param.Name, Type: param.Type}
			}
			f.Real += "_" + param.Type.Name
			f.Params = append(f.Params, param)
			if p.match(RIGHT_PAREN) {
				break
			}
			p.consume(COMMA, "Expecting ) or ,")
		}
	}
	f.Return = *p.getTypeByKind(VOID)
	if p.check(DECLARE) {
		cls, _ := p.TypeDeclare(false)
		f.Return = cls
	}
	p.consume(LEFT_BRACE, "Expecting function block")
	p.consume(EOL, "Expecting end of line")
	p.ActualFunction = f
	if p.ActualClass.Name == "" {
		p.Functions[f.Real] = f
	}
	f.Body = p.block()
	if p.ActualFunction.Returned == 0 && p.ActualFunction.Return.Kind >= STRING {
		p.error("Expecting return from funcion '"+f.Name+"'", 0)
	}
	p.EndScope()
	p.ActualFunction = Function{}
	return Node{Type: FUNCTION, Value: f}
}

func (t *Transpiler) Function(node Node) {
	f := node.Value.(Function)
	t.file.WriteString(f.Return.Real + " func_" + f.Real + "(")
	for i, p := range f.Params {
		if i > 0 {
			t.file.WriteString(", ")
		}
		t.file.WriteString(p.Type.Real + " param_" + p.Name)
	}
	t.file.WriteString(") {\n")
	for _, n := range f.Body {
		t.Transpile(n)
	}
	t.file.WriteString("}\n")
}
