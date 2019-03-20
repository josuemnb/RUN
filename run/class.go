package run

import "strings"

type Class struct {
	Name       string
	Real       string
	Super      interface{}
	Fields     map[string]Variable
	Methods    map[string]Function
	This       map[string]Function
	Kind       int
	Protection Protection
}

func (m *Module) isClass(c string) bool {
	if c, ok := m.Types[c]; ok {
		return c.Class.Name != ""
	}
	return false
}

func (m *Module) getClass(c string) Class {
	return m.Types[c].Class
}

func (t *Transpiler) Class(node Node) {
	t.insideClass = true
	class := node.Value.(Class)
	t.class = class
	t.file.WriteString("\nclass ")
	if t.Program.Type == MODULE {
		t.file.WriteString(t.Program.Name + "_")
	}
	t.file.WriteString("class_")
	t.file.WriteString(t.class.Name + " {\n")
	for n, f := range class.Fields {
		if n[0] == '_' {
			if len(n) > 2 && n[1] == '_' {
				t.file.WriteString("private:\n")
			} else {
				t.file.WriteString("protected:\n")
			}
		} else {
			t.file.WriteString("public:\n")
		}
		t.file.WriteString(f.Type.Real + " field_" + n + ";\n")
	}
	for n, m := range class.Methods {
		if n[0] == '_' {
			if len(n) > 2 && n[1] == '_' {
				t.file.WriteString("private:\n")
			} else {
				t.file.WriteString("protected:\n")
			}
		} else {
			t.file.WriteString("public:\n")
		}
		t.file.WriteString(m.Return.Real + " method_" + n + "(")
		for i, p := range m.Params {
			if i > 0 {
				t.file.WriteString(",")
			}
			t.file.WriteString(p.Type.Real + " param_" + p.Name)
		}
		t.file.WriteString(") {\n")

		for _, f := range m.Body {
			t.Transpile(f)
		}
		t.file.WriteString("}\n")
	}
	if len(class.This) > 0 {
		t.file.WriteString("public:\n")
		for _, this := range class.This {
			if t.Program.Type == MODULE {
				t.file.WriteString(t.Program.Name + "_")
			}
			t.file.WriteString("class_" + class.Name + "(")
			for i, p := range this.Params {
				if i > 0 {
					t.file.WriteString(",")
				}
				t.file.WriteString(p.Type.Real + " param_" + p.Name)
			}
			t.file.WriteString(") {\n")
			for _, f := range this.Body {
				t.Transpile(f)
			}
			t.file.WriteString("}\n")
		}
	}
	t.file.WriteString("}")
	t.class = Class{}
	t.insideClass = false
}

func (t *Transpiler) This(node Node) {
	this := node.Value.(Identifier)
	if this.Kind == FIELD {
		t.file.WriteString("field_" + this.Name)
	} else {
		t.file.WriteString("method_" + this.Name)
	}
}

func (p *Module) This() Node {
	if p.insideClass() == false || p.insideFunction() == false {
		p.error("keyword 'this' not allowed outside class function body", 0)
	}
	p.consume(DOT, "Expecting .")
	id := p.consume(IDENTIFIER, "Expecting identfier").Lexeme
	if f, ok := p.ActualClass.Fields[id]; ok {
		// println("IS FIELD", id)
		return Node{Type: THIS, Value: Identifier{Name: id, Kind: FIELD, Type: f.Type}}
	} else if !p.match(LEFT_PAREN) {
		p.error("Unknown class identifier", 1)
		params := id + p.CheckParams()
		if p.ActualClass.isMethod(params) {
			return p.GetParams(p.ActualClass.Methods[params], METHOD)
		}
		p.error("Unknown class identifier", 1)
	}
	return Node{}
}

func (p *Module) Class() Node {
	cls := p.consume(IDENTIFIER, "Expecting name")
	if p.isClass(cls.Lexeme) {
		p.error("Class already defined", 1)
	}
	var class Class
	class.Fields = make(map[string]Variable)
	class.Methods = make(map[string]Function)
	class.This = make(map[string]Function)
	// if p.Type == MODULE {
	// 	class.Name = p.Name + "_" + cls.Lexeme
	// } else {
	class.Name = cls.Lexeme
	// }
	p.consume(LEFT_BRACE, "Expecting {")
	p.consume(EOL, "Expecting end of line")
	p.ActualClass = class
	typ := Type{Name: class.Name, Class: class}
	p.addType(&typ)
	p.BeginScope()
	for !p.match(RIGHT_BRACE) {
		s := p.parse()
		if s.Type == EMPTY {
			continue
		}
		switch s.Type {
		case FUNCTION:
			f := s.Value.(Function)
			n := f.Name
			var access Protection
			if n[0] == '_' {
				if len(n) > 2 && n[1] == '_' {
					access = PRIVATE
				} else {
					access = PROTECTED
				}
			}
			f.Protection = access
			for _, fn := range f.Params {
				n += "_" + fn.Type.Name
			}
			f.Real = n
			if f.Name == "this" {
				f.Return = typ
				f.Real = class.Name
				class.This[n] = f
			} else {
				class.Methods[n] = f
			}
		case DECLARE:
			d := s.Value.(Declare)
			var p Protection
			if d.Name[0] == '_' {
				if len(d.Name) > 2 && d.Name[1] == '_' {
					p = PRIVATE
				} else {
					p = PROTECTED
				}
			}
			class.Fields[d.Name] = Variable{Name: d.Name, Type: d.Type, Protection: p}
		case EOL:
			continue
		default:
			p.error("Not allowed code", 1)
		}
	}
	// typ = Type{Name: class.Name, Class: class}
	// p.updateType(&typ)
	p.EndScope()
	p.ActualClass = Class{}
	class.Kind = typeIdx
	return Node{Type: CLASS, Value: class}
}

func (c *Class) isField(f string) bool {
	_, ok := c.Fields[f]
	return ok
}

func (c *Class) isMethod(m string) bool {
	for _, f := range c.Methods {
		if strings.HasPrefix(f.Name, m) {
			return true
		}
	}
	return false
}

func (c *Class) getFieldType(v string) Type {
	if f, ok := c.Fields[v]; ok {
		return f.Type
	}
	return Type{}
}

func (c *Class) getField(v string) Variable {
	if f, ok := c.Fields[v]; ok {
		return f
	}
	return Variable{}
}

func (c *Class) getMethodType(v string) Type {
	if f, ok := c.Methods[v]; ok {
		return f.Return
	}
	return Type{}
}
