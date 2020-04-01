package run

import (
	"fmt"
	"log"
)

func (m *Module) Parse(node *Node) {
	if !node.Parsed && node.Type != NEWLINE {
		m.parse(node)
	}
	if len(node.Children) > 0 {
		for _, n := range node.Children {
			if !n.Parsed && n.Type != NEWLINE {
				m.parse(n)
			}
		}
	}
}

func (m *Module) parse(node *Node) {
	if node.Type == KEYWORD {
		m.parseKeyword(node)
	} else if node.Type == BUILTIN {
		m.parseBuiltin(node)
	} else if node.Type == IDENTIFIER {
		m.parseIdentifier(node)
	}
}

func (m *Module) parseIdentifier(node *Node) {
	if len(node.Children) > 0 {
		if node.Children[0].Type == LEFT_PAREN {
			//test for function declaration or assigment
			if m.isFunction(node.Token.Value) {
				if m.insideFunction() == false {
					log.Fatal("Error: Function call only inside function bodys at " + fmt.Sprint(node.Token.Line))
				}
				m.parseFunctionCall(node)
				return
			} else if m.insideFunction() {
				log.Fatal("Error: Function declaration not allowed inside functions at " + fmt.Sprint(node.Token.Line))
			}
			m.FunctionDeclare(node)
			return
		} else if node.Children[0].Type == LEFT_BRACE {
			if _, _, ok := m.getName(node.Token.Value); ok {
				log.Fatal("Error: Name already in use " + fmt.Sprint(node.Token.Line))
			}
			if m.insideFunction() {
				log.Fatal("Error: Class declaration not allowed inside functions at " + fmt.Sprint(node.Token.Line))
			}
			m.parseClassDeclare(node)
		} else if node.Children[0].Type == LEFT_BRACKET {
			m.parseCollection(node)
		} else if m.insideFunction() {

		}
	} else if m.parseName(node) == false {
		m.parseNewVariable(node.Parent, node.Index)
	}
}

func (m *Module) parseName(node *Node) bool {
	if t, inter, ok := m.getName(node.Token.Value); ok {
		node.Code = node.Token.Value
		node.Parsed = true
		if node.Parent.Children[node.Index+1].Type == EQUAL {
			node.Parent.Children[node.Index+1].Code = " = "
			node.Parent.Children[node.Index+1].Parsed = true
			i, typ := m.parseExpression(node.Parent, node.Index+2)
			kind := (inter.(*Variable)).Type.Kind
			if t != VARIABLE || kind != typ {
				log.Fatal("Error: Mismatches types at " + fmt.Sprint(node.Token.Line))
			}
			if node.Parent.Children[i].Type != NEWLINE {
				log.Fatal("Error: Expecting end of line at " + fmt.Sprint(node.Token.Line))
			}
			node.Parent.Children[i].Code = ";\n"
		}
		return true
	}
	return false
}

func (m *Module) parseKeyword(node *Node) {
	switch node.Token.Value {
	case "main":
		m.parseMain(node)
	case "loop":
		m.parseLoop(node)
	case "return":
		m.parseReturn(node)
	case "this":
		m.parseThis(node)
	}
}

func (m *Module) parseBuiltin(node *Node) {
	switch node.Token.Value {
	case "print":
		m.parsePrint(node, false)
	case "println":
		m.parsePrint(node, true)
	}
}

func intToType(t int) (s string) {
	switch t {
	default:
		s = ""
	}
	return
}

func stringToType(v string) (s string) {
	switch v {
	case "number", "bool", "real", "string", "byte":
		s = v
	}
	return
}

func parseParams(params *Node) {
	l := len(params.Children)
	index := 0
	params.Parsed = true
	for index < l {
		n := params.Children[index]
		if n.Type != IDENTIFIER {
			log.Fatal("Error: Expecting a Identifier at " + fmt.Sprint(params.Token.Line))
		}
		index++
		name := n.Token.Value
		if l < index {
			n.Code = "Run_string " + name
			break
		}
		n1 := params.Children[index]
		index++
		if n1.Type == COMMA {
			n.Code = "Run_string " + name
			n1.Code = ", "
			continue
		} else if n1.Type != DECLARE {
			log.Fatal("Error: Expecting a comma or ':' at " + fmt.Sprint(params.Token.Line))
		}
		if l < index {
			log.Fatal("Error: Expecting a Type at " + fmt.Sprint(params.Token.Line))
		}
		n2 := params.Children[index]
		if n2.Type != IDENTIFIER {
			log.Fatal("Error: Expecting a Type at " + fmt.Sprint(params.Token.Line))
		}
		n2.Code = name
		n1.Code = " "
		n.Code = stringToType(n2.Token.Value)
		index++
		if l < index {
			log.Fatal("Error: Expecting a ')' or comma at " + fmt.Sprint(params.Token.Line))
		}
		if params.Children[index].Type == RIGHT_PAREN {
			params.Children[index].Code = ") "
			break
		}
		if params.Children[index].Type != COMMA {
			log.Fatal("Error: Expecting a ')' or comma at " + fmt.Sprint(params.Token.Line))
		}
		index++
	}
}

func (m *Module) parseBody(body *Node) {
	l := len(body.Children)
	body.Parsed = true
	for i := 0; i < l; i++ {
		m.Parse(body.Children[i])
	}
}
