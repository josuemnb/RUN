package run

import (
	"os"
	"strings"
)

type Class struct {
	Name       string
	Real       string
	Super      []Class
	Fields     map[string]Variable
	Methods    map[string]Function
	Operators  map[string]Operator
	Virtuals   []string
	This       map[string]Function
	Kind       int
	Cpp        []string
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

type Operator struct {
	op     string
	params []Param
	Return Type
	Body   []Node
	Real   string
}

func (m *Module) operator() Node {
	var params []Param
	op := m.previous().Lexeme
	var ret Type
	var real string
	m.BeginScope()
	if m.match(LEFT_PAREN) {
		params, real = m.ParseParams()
	}
	if m.check(DECLARE) {
		ret, _ = m.TypeDeclare(false, true)
	}
	m.consume(LEFT_BRACE, "Expecting function block")
	m.consume(EOL, "Expecting end of line")
	m.ActualFunction = Function{Name: op, Params: params, Return: ret}
	Body := m.block()
	if m.ActualFunction.Returned == 0 && m.ActualFunction.Return.Kind >= STRING {
		m.error("Expecting return from funcion '"+op+"'", 0)
	}
	m.ActualFunction = Function{}
	m.EndScope()
	return Node{Type: OPERATOR, Value: Operator{op: op, params: params, Return: ret, Body: Body, Real: real}}
}

func (t *Transpiler) Class(node Node) {
	t.insideClass = true
	class := node.Value.(Class)
	t.class = class
	file, err := os.OpenFile(class.Name+".h", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	temp := t.file
	t.file = file
	t.file.WriteString("#pragma once\n\n#include \"run/libc/run.h\"\n\nclass ")
	if t.Program.Type == MODULE {
		t.file.WriteString(t.Program.Name + "_")
	}
	t.file.WriteString("class_" + t.class.Name)
	if len(class.Super) > 0 {
		t.file.WriteString(": ")
		for i, s := range class.Super {
			if i > 0 {
				t.file.WriteString(", ")
			}
			t.file.WriteString("class_" + s.Name)
		}
	}
	t.file.WriteString(" {\n")
	if len(class.Cpp) > 0 {
		t.file.WriteString("private:\n")
		for _, cpp := range class.Cpp {
			t.file.WriteString(cpp + "\n")
		}
	}
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
		if f.Type.IsInterface {
			t.file.WriteString("static ")
			t.InterfaceDeclare(Declare{Name: n, Type: f.Type})
		} else {
			t.file.WriteString(f.Type.Real + " field_" + n + ";\n")
		}
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
		if m.IsVirtual {
			t.file.WriteString("virtual ")
		}
		t.file.WriteString(m.Return.Real + " method_" + n + "(")
		for i, p := range m.Params {
			if i > 0 {
				t.file.WriteString(",")
			}
			if p.Type.IsInterface {
				t.file.WriteString(p.Type.Real + " param_inter_" + p.Name)
			} else {
				t.file.WriteString(p.Type.Real + " param_" + p.Name)
			}
		}
		t.file.WriteString(")")
		if m.IsVirtual {
			t.file.WriteString(";\n")
		} else {
			t.file.WriteString(" {\n")

			for _, f := range m.Body {
				t.Transpile(f)
			}
			t.file.WriteString("}\n")
		}
	}
	if len(class.Operators) > 0 {
		t.file.WriteString("public:\n")
		for _, op := range class.Operators {
			if op.Return.Name == class.Name {
				t.file.WriteString(op.Return.Real + " &operator" + op.op + "(")
			} else {
				t.file.WriteString(op.Return.Real + " operator" + op.op + "(")
			}
			for i, p := range op.params {
				if i > 0 {
					t.file.WriteString(",")
				}
				if p.Type.IsInterface {
					t.file.WriteString(p.Type.Real + " param_inter_" + p.Name)
				} else if p.Type.Kind >= STRING {
					t.file.WriteString(p.Type.Real + " &param_" + p.Name)
				} else {
					t.file.WriteString(p.Type.Real + " param_" + p.Name)
				}
			}
			t.file.WriteString(") {\n")

			for _, f := range op.Body {
				t.Transpile(f)
			}
			t.file.WriteString("}\n")
		}
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
	t.file.WriteString("};")
	t.class = Class{}
	t.insideClass = false
	t.file.Close()
	t.file = temp
}

func (t *Transpiler) This(node Node) {
	this := node.Value.(Identifier)
	if this.Kind == FIELD {
		t.file.WriteString("field_" + this.Name)
	} else if this.Kind == THIS {
		t.file.WriteString("*this")
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

func (p *Module) Class(extends bool) Node {
	cls := p.consume(IDENTIFIER, "Expecting name")
	if p.isClass(cls.Lexeme) {
		p.error("Class already defined", 1)
	}
	var class Class
	class.Fields = make(map[string]Variable)
	class.Methods = make(map[string]Function)
	class.This = make(map[string]Function)
	class.Operators = make(map[string]Operator)
	class.Cpp = make([]string, 0)
	class.Virtuals = make([]string, 0)
	class.Super = make([]Class, 0)

	class.Name = cls.Lexeme
	if extends {
		p.consume(EXTENDS, "Expecting <- extends symbol")
		for {
			s := p.consume(IDENTIFIER, "Expecting name")
			if !p.isClass(s.Lexeme) {
				p.error("Class undefined", 1)
			}
			class.Super = append(class.Super, p.getClass(s.Lexeme))
			if p.check(LEFT_BRACE) {
				break
			}
			p.consume(COMMA, "Expecting comma")
		}
	}
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
		case OPERATOR:
			o := s.Value.(Operator)
			if _, ok := class.Operators[o.op+o.Real]; ok {
				p.error("Operator already defined", 0)
			}
			class.Operators[o.op+o.Real] = o
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
				if f.IsVirtual {
					class.Virtuals = append(class.Virtuals, n)
				}
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
		case CPP:
			cpp := s.Value.(Cpp)
			class.Cpp = append(class.Cpp, cpp.Body)
		default:
			p.error("Not allowed code", 1)
		}
	}
	// class.Kind = typeIdx
	// typ = Type{Name: class.Name, Class: class}
	// p.updateType(&Type{Name: class.Name, Class: class})
	p.EndScope()
	p.ActualClass = Class{}
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
