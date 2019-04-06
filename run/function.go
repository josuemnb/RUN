package run

type Function struct {
	Params     []Param
	Body       []Node
	Return     Type
	Name       string
	Protection Protection
	Returned   int
	Real       string
	IsVirtual  bool
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

func (p *Module) PointerFunction() Node {
	var f Function
	t := p.advance()
	if p.isFunction(t.Lexeme) {
		p.error("Function name already assigned", 1)
	}
	if _, ok := p.isVar(t.Lexeme); ok {
		p.error("Name already assigned to var", 1)
	}
	f.Name = t.Lexeme
	if p.Type == MODULE {
		f.Real = p.Name + "_"
	}
	// if block == false {
	// 	f.Real += "param_"
	// }
	f.Real += f.Name
	p.advance()
	var real string
	f.Params, real = p.ParseParams()
	f.Real += real
	f.Return = *p.getTypeByKind(VOID)
	if p.ActualClass.Name == "" {
		p.Functions[f.Real] = f
	}
	return Node{Type: FUNCTION, Value: f}
}

func (p *Module) ParseParams() (params []Param, real string) {
	params = make([]Param, 0)
	if !p.match(RIGHT_PAREN) {
		for {
			var param Param
			n := p.advance()
			for _, pm := range params {
				if pm.Name == n.Lexeme {
					p.error("Name already assigned", 0)
				}
			}
			param.Name = n.Lexeme
			if p.check(DECLARE) {
				cls, array := p.TypeDeclare(true, true)
				param.Type = cls
				p.Scopes[p.CurScope][param.Name] = Variable{Name: param.Name, Type: cls, Array: array}
				// } else if p.check(LEFT_PAREN) {
				// 	p.rollBack()
				// 	pfunc := p.PointerFunction()
				// 	pValue := pfunc.Value.(Function)
				// 	param.Type = Type{Kind: FUNCTION, Function: pfunc, IsFunction: true}
				// 	param.Type.Name = "ptr_" + pValue.Real
				// 	p.Scopes[p.CurScope][param.Name] = Variable{Name: param.Name, Type: param.Type}
			} else {
				param.Type = *p.getTypeByName("string")
				p.Scopes[p.CurScope][param.Name] = Variable{Name: param.Name, Type: param.Type}
			}
			real += "_" + param.Type.Name
			params = append(params, param)
			if p.match(RIGHT_PAREN) {
				break
			}
			p.consume(COMMA, "Expecting ) or ,")
		}
	}
	return
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
	// if block == false {
	// 	f.Real += "param_"
	// }
	f.Real += f.Name
	p.advance()
	p.BeginScope()
	var real string
	f.Params, real = p.ParseParams()
	f.Real += real
	f.Return = *p.getTypeByKind(VOID)
	if p.check(DECLARE) {
		cls, _ := p.TypeDeclare(false, true)
		f.Return = cls
	}
	if p.ActualClass.Name == "" {
		p.Functions[f.Real] = f
	}
	// p.consume(LEFT_BRACE, "Expecting function block")
	if p.match(LEFT_BRACE) {
		p.consume(EOL, "Expecting end of line")
		p.ActualFunction = f
		f.Body = p.block()
		if p.ActualFunction.Returned == 0 && p.ActualFunction.Return.Kind >= STRING {
			p.error("Expecting return from funcion '"+f.Name+"'", 0)
		}
		p.ActualFunction = Function{}
	} else if p.insideClass() {
		f.IsVirtual = true
	} else {
		p.error("Not allowed", 1)
	}
	p.EndScope()
	return Node{Type: FUNCTION, Value: f}
}

func (t *Transpiler) PointerFunction(node Node) {
	f := node.Value.(Function)

	t.file.WriteString(f.Return.Real + " (*func_" + f.Real + ")(")
	for i, p := range f.Params {
		if i > 0 {
			t.file.WriteString(", ")
		}
		// if p.Type.IsFunction {
		// 	// t.file.WriteString(p.Type.Real + " func_" + p.Name)
		// 	t.PointerFunction(p.Type.Function)
		// } else {
		t.file.WriteString(p.Type.Real + " param_" + p.Name)
		// }
	}
	t.file.WriteString(")")
}

func (t *Transpiler) Function(node Node) {
	f := node.Value.(Function)
	t.file.WriteString(f.Return.Real + " func_" + f.Real + "(")
	for i, p := range f.Params {
		if i > 0 {
			t.file.WriteString(", ")
		}
		if p.Type.IsInterface {
			t.file.WriteString(p.Type.Real + " param_inter_" + p.Name)
		} else {
			t.file.WriteString(p.Type.Real + " param_" + p.Name)
		}
	}
	t.file.WriteString(") {\n")
	for _, n := range f.Body {
		t.Transpile(n)
	}
	t.file.WriteString("}\n")
}
