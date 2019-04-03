package run

type Interface struct {
	Name      string
	Functions map[string]Function
	Init      bool
}

func (p *Module) InterfaceDeclare(in Interface, n string) (cls Type) {
	// cls = Type{Name: in.Name}
	interf := Interface{Name: in.Name, Functions: in.Functions, Init: true}
	cls.Name = in.Name
	p.consume(LEFT_BRACE, "Expecting {")
	p.consume(EOL, "Expecting end of line")
	count := 0
	for !p.match(RIGHT_BRACE) {
		f := p.consume(IDENTIFIER, "Expecting interface function name")
		p.consume(LEFT_PAREN, "Expecting (")
		p.BeginScope()
		params, real := p.ParseParams()
		fi, ok := interf.Functions[f.Lexeme+real]
		if !ok {
			p.error("Unknown Function or wrong parameters for interface", 1)
		}
		p.consume(LEFT_BRACE, "Expeting {")
		fi.Params = params
		p.ActualFunction = fi
		fi.Name = n + real
		fi.Real = f.Lexeme + real
		fi.Body = p.block()
		p.ActualFunction = Function{}
		p.EndScope()
		fi.Return = *p.getTypeByKind(VOID)
		interf.Functions[fi.Real] = fi
		count++
	}
	if count != len(in.Functions) {
		p.error("Functions implementation different than declared", 0)
	}
	cls.IsInterface = true
	cls.Interface = interf
	// cls.Kind = p.getTypeByName(interf.Name).Kind
	return
}

func (p *Module) Interface() Node {
	n := p.consume(IDENTIFIER, "Expecting name for interface")
	if _, ok := p.Types[n.Lexeme]; ok {
		p.error("Name used as as type", 0)
	}
	p.advance()
	p.consume(LEFT_BRACE, "Expecting block for interface")
	p.consume(EOL, "Expecting new line")
	functions := make(map[string]Function)
	for !p.match(RIGHT_BRACE) {
		f := p.consume(IDENTIFIER, "Expecting name of function")
		p.consume(LEFT_PAREN, "Expecting (")
		params, real := p.ParseParams()
		real = f.Lexeme + real
		functions[real] = Function{Params: params, Real: real, Name: f.Lexeme}
		p.consume(EOL, "Expecting end of line")
	}
	// if p.Type == MODULE {
	// 	n.Lexeme = p.Name + "_" + n.Lexeme
	// }
	interf := Interface{Name: n.Lexeme, Functions: functions}
	// p.Interfaces[n.Lexeme] = interf
	p.addType(&Type{Name: n.Lexeme, IsInterface: true, Interface: interf})
	return Node{Type: INTERFACE, Value: interf}
}

func (t *Transpiler) Interface(node Node) {
	i := node.Value.(Interface)
	t.file.WriteString("typedef struct ")
	if t.Program.Type == MODULE {
		t.file.WriteString(t.Program.Name + "_")
	}
	t.file.WriteString(i.Name + " {\n")
	for _, f := range i.Functions {
		t.file.WriteString("void (*" + f.Real + ")(")
		for c, p := range f.Params {
			if c > 0 {
				t.file.WriteString(",")
			}
			t.file.WriteString(p.Type.Real + " " + p.Name)
		}
		t.file.WriteString(");\n")
	}
	t.file.WriteString("}")
}

func (t *Transpiler) InterfaceDeclare(d Declare) {
	if d.Type.Interface.Init == false {
		t.file.WriteString(d.Type.Interface.Name + " inter_" + d.Name + ";\n")
		return
	}
	for _, f := range d.Type.Interface.Functions {
		t.file.WriteString("void inter_" + d.Name + "_" + f.Real + "(")
		for i, p := range f.Params {
			if i > 0 {
				t.file.WriteString(",")
			}
			t.file.WriteString(p.Type.Real + " param_" + p.Name)
		}
		if d.Type.Interface.Init {
			t.file.WriteString(") {\n")
			for _, b := range f.Body {
				t.Transpile(b)
			}
			t.file.WriteString("}\n")
		} else {
			t.file.WriteString(");\n")
		}
	}
	t.file.WriteString(t.Program.getTypeByName(d.Type.Interface.Name).Name + " inter_" + d.Name + " = {\n")
	i := len(d.Type.Interface.Functions)
	for _, f := range d.Type.Interface.Functions {
		t.file.WriteString("." + f.Real + "= inter_" + d.Name + "_" + f.Real)
		if i > 1 {
			t.file.WriteString(",")
		}
		t.file.WriteString("\n")
		i--
	}
	t.file.WriteString("};\n")
}
