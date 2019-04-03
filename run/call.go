package run

import (
	"bytes"
)

type Call struct {
	Name   string
	Args   []Node
	Kind   int
	Return Type
}

type Get struct {
	Name  string
	Kind  int
	Base  Node
	Value Node
}

func (p *Module) Instantiate(n string) Node {
	class := p.getClass(n)
	f := "this" + p.CheckParams()
	// p.CurToken++
	if fun, ok := class.This[f]; ok {
		return p.GetParams(fun, THIS)
	}
	return Node{}
}

func (p *Module) CheckParams() string {
	current := p.CurToken
	if p.check(RIGHT_PAREN) {
		// p.CurToken--
		return ""
	}
	var buff bytes.Buffer
	for {
		arg := p.assignment()
		t := p.typeOf(arg)
		if t.Kind == 0 && arg.Type == IDENTIFIER {
			id := arg.Value.(Identifier)
			// println(id.Name)
			if p.isFunction(id.Name) {
				t.Name = "ptr_func"
			}
		}
		if t.Kind == QUOTE {
			t = *p.getTypeByKind(STRING)
		}
		buff.WriteString("_" + t.Name)
		if p.match(RIGHT_PAREN) {
			break
		}
		p.consume(COMMA, "Expecting comma or )")
	}
	p.CurToken = current
	return buff.String()
}

func (p *Module) GetParams(f Function, kind int) Node {
	if len(f.Params) == 0 {
		p.consume(RIGHT_PAREN, "Expecting )")
		return Node{Type: CALL, Value: Call{Name: f.Real, Kind: kind, Args: nil, Return: f.Return}}
	}
	args := make([]Node, 0)
	for i, param := range f.Params {
		if i > 0 {
			p.consume(COMMA, "Expecting comma")
		}
		arg := p.assignment()
		t := p.typeOf(arg)
		if t.Kind == 0 && arg.Type == IDENTIFIER {
			id := arg.Value.(Identifier)
			// println(id.Name)
			if p.isFunction(id.Name) {
				t.Kind = FUNCTION
				id.Kind = FUNCTION
				arg.Value = id
			}
		}
		if t.Kind == QUOTE {
			t = *p.getTypeByKind(STRING)
		}
		if t.Kind != param.Type.Kind {
			p.error("Mismatch kinds", 1)
		}
		args = append(args, arg)
	}
	p.consume(RIGHT_PAREN, "Expecting )")
	return Node{Type: CALL, Value: Call{Name: f.Real, Kind: kind, Args: args, Return: f.Return}}
}

func (m *Module) parseFunction(n string) Node {
	params := n + m.CheckParams()
	if m.insideClass() && m.ActualClass.isMethod(params) {
		return m.GetParams(m.ActualClass.Methods[params], METHOD)
	} else if m.isFunction(params) {
		return m.GetParams(m.Functions[params], FUNCTION)
	} else {
		m.error("Unknwon identifier", 1)
	}
	return Node{}
}

func (p *Module) Call(parent interface{}) Node {
	op := "."
	e := p.primary()
	if e.Type == IDENTIFIER {
		id := e.Value.(Identifier)
		if p.match(LEFT_PAREN) {
			if parent != nil {
				switch parent.(type) {
				case Module:
					m := parent.(Module)
					n := m.Name + "_" + id.Name + p.CheckParams()
					if m.isFunction(n) {
						f := m.Functions[n]
						e = p.GetParams(f, FUNCTION)
						if f.Return.IsInterface {
							parent = f.Return.Interface
						} else if f.Return.Kind >= STRING {
							parent = f.Return.Class
						} else {
							return e
						}
					}
				case Interface:
					i := parent.(Interface)
					n := id.Name + p.CheckParams()
					if f, ok := i.Functions[n]; ok {
						e = p.GetParams(f, INTERFACE)
						if f.Return.IsInterface {
							parent = f.Return.Interface
						} else if f.Return.Kind >= STRING {
							parent = f.Return.Class
						} else {
							return e
						}
					} else {
						p.error("Unknown interface function", 1)
					}
				case Class:
					c := parent.(Class)
					n := id.Name + p.CheckParams()
					if f, ok := c.Methods[n]; ok {
						if f.Name != id.Name {
							p.error("Unknown identifier", 1)
						}
						e = p.GetParams(f, METHOD)
						if f.Return.IsInterface {
							parent = f.Return.Interface
						} else if f.Return.Kind >= STRING {
							parent = f.Return.Class
						} else {
							return e
						}
						// } else if f, ok := c.Fields[id.Name]; ok {
						// 	right = Node{Type: IDENTIFIER, Value: Identifier{Name: id.Name, Kind: FIELD, Type: f.Type}}
					} else {
						p.error("Unknown Identifier", 1)
					}
				case Function:
					f := parent.(Function)
					if f.isParam(id.Name) {
						p.error("Name already assigned", 0)
					}
					e = p.GetParams(f, FUNCTION)
					if f.Return.IsInterface {
						parent = f.Return.Interface
					} else if f.Return.Kind >= STRING {
						parent = f.Return.Class
					} else {
						return e
					}
				}
			} else if p.isClass(id.Name) {
				e = p.Instantiate(id.Name)
			} else {
				e = p.parseFunction(id.Name)
				tp := p.typeOf(e)
				if tp.IsInterface {
					parent = tp.Interface
				} else if tp.Kind >= STRING {
					parent = tp.Class
				} else {
					return e
				}
			}
		} else if p.insideFunction() && p.ActualFunction.isParam(id.Name) {
			v := p.ActualFunction.getParam(id.Name)
			if v.Type.IsInterface {
				id.Name = "inter_" + id.Name
			}
			e = Node{Type: IDENTIFIER, Value: Identifier{Name: id.Name, Kind: PARAM, Type: v.Type}}
			if v.Type.IsInterface {
				// op = "_"
				parent = v.Type.Interface
			} else if v.Type.Kind >= STRING {
				parent = v.Type.Class
			} else {
				return e
			}
		} else if p.insideClass() && p.ActualClass.isField(id.Name) {
			v := p.ActualClass.Fields[id.Name]
			e = Node{Type: IDENTIFIER, Value: Identifier{Name: id.Name, Kind: FIELD, Type: v.Type}}
			if v.Type.IsInterface {
				parent = v.Type.Interface
			} else if v.Type.Kind >= STRING {
				parent = v.Type.Class
			} else {
				return e
			}
		} else if v, ok := p.isVar(id.Name); ok {
			e = Node{Type: IDENTIFIER, Value: Identifier{Name: id.Name, Kind: VARIABLE, Type: v.Type}}
			if v.Type.IsInterface {
				// op = "_"
				parent = v.Type.Interface
			} else if v.Type.Kind >= STRING {
				parent = v.Type.Class
			} else {
				return e
			}
		} else if m, ok := p.Modules[id.Name]; ok {
			parent = *m
			op = ""
		} else if f, ok := p.Functions[id.Name]; ok {
			println("FUNCTION", f.Name)
		}
		if p.match(DOT) && parent != nil {
			return Node{Type: BINARY, Value: Binary{e, p.Call(parent), op}}
		}
	}
	return e
}

func (t *Transpiler) Call(node Node) {
	c := node.Value.(Call)
	if c.Kind == FUNCTION {
		t.file.WriteString("func_" + c.Name + "(")
	} else if c.Kind == METHOD {
		t.file.WriteString("method_" + c.Name + "(")
	} else if c.Kind == CLASS {
		t.file.WriteString("class_" + c.Name + "(")
	} else if c.Kind == THIS {
		t.file.WriteString("(")
	} else if c.Kind == INTERFACE {
		t.file.WriteString(c.Name + "(")
	} else {

	}
	if c.Args != nil {
		for i, p := range c.Args {
			if i > 0 {
				t.file.WriteString(", ")
			}
			t.Transpile(p)
		}
	}
	t.file.WriteString(")")
}
