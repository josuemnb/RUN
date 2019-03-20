package run

import (
	"strconv"
)

type Declare struct {
	Name  string
	Array int
	Type  Type
}

func (p *Module) TypeDeclare(inst bool) (cls Type, array int) {
	p.consume(DECLARE, "expecting :")
	t := p.advance()
	if t.Lexeme == "[" {
		t = p.advance()
		if t.Type == NUMBER {
			array, _ = strconv.Atoi(t.Lexeme)
			t = p.advance()
		}
		if t.Lexeme != "]" {
			p.error("Expecting ]", 0)
		}
		if array == 0 {
			array = -1
			cls = p.List()
		} else {
			cls = p.Array(array)
		}
	} else if t.Lexeme == "map" {
		cls = p.Map()
	} else if t.Lexeme == "list" {
		cls = p.List()
	} else if t.Lexeme == "stack" {

	} else if t.Lexeme == "array" {
		cls = p.Array(-1)
	} else if m, ok := p.Modules[t.Lexeme]; ok {
		p.consume(DOT, "Expecting .")
		t = p.consume(IDENTIFIER, "Expecting Identfier")
		typ := m.getTypeByName(t.Lexeme)
		if typ == nil || typ.Name == "" {
			p.error("Unknown identifier", 0)
		}
		cls = *typ
	} else {
		cls = *p.getTypeByName(t.Lexeme)
		if cls.Name == "" {
			if p.ActualClass.Name == t.Lexeme {
				cls = Type{Name: p.ActualClass.Name, Kind: cls.Kind, Class: p.ActualClass}
			} else {
				p.error("unknown type: "+t.Lexeme, 0)
			}
		}
		if inst && cls.Kind >= STRING && len(cls.Class.This) > 0 {
			ok := false
			for _, this := range cls.Class.This {
				if len(this.Params) == 0 {
					ok = true
					break
				}
			}
			if !ok {
				p.error("Class doesn't allow simple declaration", 1)
			}
		}
	}
	return
}

func (p *Module) Declare() Node {
	n := p.advance()
	if _, ok := p.isVar(n.Lexeme); ok {
		p.error("Variable name in use", 0)
	}
	cls, array := p.TypeDeclare(true)
	p.consume(EOL, "expecting end of line")
	p.Scopes[p.CurScope][n.Lexeme] = Variable{Name: n.Lexeme, Type: cls, Array: array}
	return Node{Type: DECLARE, Value: Declare{Name: n.Lexeme, Type: cls, Array: array}}
}

func (t *Transpiler) Declare(node Node) {
	d := node.Value.(Declare)
	t.file.WriteString(d.Type.Real + " var_" + d.Name)
	t.file.WriteString(";\n")
}
