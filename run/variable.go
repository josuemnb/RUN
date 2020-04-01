package run

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Variable struct {
	Name    string
	Address string
	Type    *Type
	Access  int
	Array   int
	Parent  interface{}
}

func (v *Variable) isArray() bool {
	return v.Array > -1
}

func (m *Module) parseNewVariable(node *Node, index int) {
	l := len(node.Children)
	name := node.Children[index].Token.Value
	node.Children[index].Parsed = true
	array := -1
	if l > index+2 {
		var typ *Type
		if node.Children[index+1].Type == EQUAL {
			node.Children[index+1].Parsed = true
			node.Children[index+1].Code = " = "
			if node.Children[index+2].Token.Value == "new" {
				n := node.Children[index+2]
				n.Parsed = true
				if len(n.Children) == 0 || n.Children[0].Type != LEFT_PAREN || n.Children[1].Type != RIGHT_PAREN {
					log.Fatal("Error: Badformed line at " + fmt.Sprint(n.Token.Line))
				}
				n.Children[0].Parsed = true
				n.Children[1].Parsed = true
				childs := n.Children[0].Children
				l := len(childs)
				if l != 3 {
					log.Fatal("Error: Badformed line at " + fmt.Sprint(n.Token.Line))
				}
				if childs[0].Type != IDENTIFIER || childs[1].Type != COMMA || childs[2].Type != NUMBER {
					log.Fatal("Error: Badformed line at " + fmt.Sprint(n.Token.Line))
				}
				typ = m.getTypeByName(childs[0].Token.Value)
				node.Children[index].Code = typ.Real + " *" + name
				childs[0].Parsed = true
				childs[1].Parsed = true
				childs[2].Parsed = true
				array, _ = strconv.Atoi(childs[2].Token.Value)
				if array == 0 {
					log.Fatal("Error: Array size must be bigger than 0 at " + fmt.Sprint(n.Token.Line))
				}
				childs[2].Code = " gcnew(" + childs[0].Token.Value + "," + childs[2].Token.Value + ");\n"
				goto end
			}
			i, t := m.parseExpression(node, index+2)
			if len(node.Children) <= i || node.Children[i].Type != NEWLINE {
				log.Fatal("Error: Expecting a new line at " + fmt.Sprint(node.Token.Line))
			}
			if t == -1 {
				log.Fatal("Error: Type is missing at " + fmt.Sprint(node.Token.Line))
			}
			node.Children[i].Parsed = true
			node.Children[i-1].Code += ";\n"
			typ = m.getTypeByIndex(t)
			if typ == nil {
				log.Fatal("Error: Unknown type name at " + fmt.Sprint(node.Token.Line))
			}
		} else if node.Children[index+1].Type == DECLARE {
			node.Children[index+1].Parsed = true
			n := node.Children[index+2]
			n.Code += ";\n"
			if n.Type == LEFT_BRACKET {
				n.Parsed = true
				if node.Children[index+3].Type != RIGHT_BRACKET {
					log.Fatal("Error: Expecting a ']' " + fmt.Sprint(node.Token.Line))
				}
				node.Children[index+3].Parsed = true
				array = 0
				if len(n.Children) > 0 {
					if n.Children[0].Type != NUMBER || len(n.Children) != 1 {
						log.Fatal("Error: Expecting a number " + fmt.Sprint(node.Token.Line))
					}
					n.Children[0].Parsed = true
					array, _ = strconv.Atoi(n.Children[0].Token.Value)
				}
				n = node.Children[index+4]
			}
			if n.Type != IDENTIFIER {
				log.Fatal("Error: Expecting a type " + fmt.Sprint(node.Token.Line))
			}
			n.Parsed = true
			typ = m.getTypeByName(n.Token.Value)
			if typ == nil {
				log.Fatal("Error: Unknown type name at " + fmt.Sprint(node.Token.Line))
			}
		} else {
			log.Fatal("Error: Badformmed line at " + fmt.Sprint(node.Token.Line))
		}
		node.Children[index].Code = typ.Real + " " + name
		if array > -1 {
			if array == 0 {
				node.Children[index].Code = typ.Real + " *" + name
			} else {
				node.Children[index].Code += "["
				if array > 0 {
					node.Children[index].Code += fmt.Sprint(array)
				}
				node.Children[index].Code += "]"
			}
		}
	end:
		v := &Variable{Name: name, Type: typ, Access: PUBLIC, Array: array}
		if m.insideFunction() {
			v.Parent = m.CurrentFunction
			m.CurrentFunction.Variables[name] = v
		} else if m.insideClass() {
			if strings.HasPrefix(name, "_") {
				if strings.HasPrefix(name, "__") {
					node.Children[index].Code = Values.PRIVATE + node.Children[index].Code
					v.Access = PRIVATE
				} else {
					node.Children[index].Code = Values.PROTECTED + node.Children[index].Code
					v.Access = PROTECTED
				}
			} else {
				node.Children[index].Code = Values.PUBLIC + node.Children[index].Code
				v.Access = PUBLIC
			}
			v.Parent = m.CurrentClass
			m.CurrentClass.Fields[name] = v
		} else {
			v.Parent = m
			m.Variables[name] = v
		}
	} else {
		log.Fatal("Error: Expecting an assignment at " + fmt.Sprint(node.Token.Line))
	}
}
