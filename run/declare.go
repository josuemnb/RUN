package run

import (
	"strconv"
)

type Declare struct {
	Name        string
	Array       int
	Type        Type
	Instantiate Node
	IsVolatile  bool
}

func (p *Module) TypeDeclare(inst bool, inFunction bool) (cls Type, array int) {
	original := p.previous().Lexeme
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
		if cls.IsInterface {
			if inFunction == false && p.check(LEFT_BRACE) {
				cls = p.InterfaceDeclare(cls.Interface, original)
				cls.Kind = typ.Kind
				return
			}
		}
	} else {
		c := p.getTypeByName(t.Lexeme)
		if c == nil {
			p.error("Type not found", 1)
		}
		if c.IsInterface {
			if inFunction == false && p.check(LEFT_BRACE) {
				cls = p.InterfaceDeclare(c.Interface, original)
				cls.Kind = c.Kind
				return
			}
		}
		cls = *c
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
	var inst Node
	n := p.advance()
	if _, ok := p.isVar(n.Lexeme); ok {
		p.error("Variable name in use", 0)
	}
	cls, array := p.TypeDeclare(true, p.insideFunction())
	if cls.Kind > STRING && cls.IsInterface == false {
		if len(cls.Class.This) > 0 {
			p.consume(LEFT_PAREN, "Expecting (")
			f := "this" + p.CheckParams()
			// p.CurToken++
			if fun, ok := cls.Class.This[f]; ok {
				inst = p.GetParams(fun, THIS)
			}
		}
	}
	p.consume(EOL, "expecting end of line")
	vol := false
	if p.CurScope == 0 {
		vol = true
	}
	p.Scopes[p.CurScope][n.Lexeme] = Variable{Name: n.Lexeme, Type: cls, Array: array}
	return Node{Type: DECLARE, Value: Declare{Name: n.Lexeme, Type: cls, Array: array, Instantiate: inst, IsVolatile: vol}}
}

func (t *Transpiler) Declare(node Node) {
	d := node.Value.(Declare)
	if d.Type.IsInterface {
		t.InterfaceDeclare(d)
	} else {
		if d.IsVolatile {
			t.file.WriteString("volatile ")
		}
		t.file.WriteString(d.Type.Real + " var_" + d.Name)
		if d.Instantiate.Type != 0 {
			// t.file.WriteString("(")
			t.Transpile(d.Instantiate)
			// t.file.WriteString(")")
		}
		t.file.WriteString(";\n")
	}
}
