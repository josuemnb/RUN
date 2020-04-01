package run

import (
	"fmt"
	"log"
	"os"
)

func (m *Module) parseLoop(node *Node) {
	node.Parsed = true
	l := len(node.Children)
	if l < 2 || node.Children[l-2].Type != LEFT_BRACE || node.Children[l-1].Type != RIGHT_BRACE {
		log.Fatal("Error: Badformed loop at " + fmt.Sprint(node.Token.Line))
		os.Exit(-1)
	}
	if l == 2 {
		node.Code = "while(1) {\n"
	} else {
		index := 0
		node.Code = "for("
		n := node.Children[index]
		n.Parsed = true
		name := "loop_" + fmt.Sprint(node.Token.Line) + "_" + fmt.Sprint(node.Token.Col)
		if n.Type == IDENTIFIER {
			name = n.Token.Value
			if v := m.getVariable(name); v != nil {
				if v.Type.Kind != NUMBER && v.Type.Kind != REAL {
					log.Fatal("Error: Only number and real values allowed in loop at " + fmt.Sprint(node.Token.Line))
					os.Exit(-1)
				}
				n.Code += name
			} else {
				n.Code += "int " + name
			}
			index++
			n = node.Children[index]
		} else {
			node.Code += "int " + name
		}
		n.Parsed = true
		if n.Type == EQUAL {
			n.Code += " = "
			index = m.parseLoopValue(node, index+1)
			node.Children[index-1].Code += ";"
			n = node.Children[index]
		} else if n.Type != RANGE {
			index = m.parseLoopValue(node, index)
			n.Code = " = " + n.Code
			n = node.Children[index]
			n.Code += ";"
		} else {
			n.Code += " = 0;"
		}
		n.Parsed = true
		if n.Type == RANGE {
			n.Code += name + "<"
			index = m.parseLoopValue(node, index+1)
			node.Children[index-1].Code += ";"
			n = node.Children[index]
		} else {
			n.Code += ";"
		}
		n.Parsed = true
		n.Code += name
		if n.Type == COMMA {
			if node.Children[index+1].Type == MINUS {
				index++
				n.Code += "-="
			} else {
				n.Code += "+="
			}
			index = m.parseLoopValue(node, index+1)
			node.Children[index].Code = ") {\n"
		} else {
			n.Code += "++) {"
		}
	}
	idx := l - 1
	if node.Children[idx].Type == NEWLINE {
		idx--
	}
	idx--
	m.parseBody(node.Children[idx])
	node.Children[idx+1].Code = "\n}\n"
}

func (m *Module) parseLoopValue(node *Node, index int) int {
	temp, t := m.parseExpression(node, index)
	if t != NUMBER && t != REAL {
		log.Fatal("Error: Only number and real values allowed in loop at " + fmt.Sprint(node.Token.Line))
		os.Exit(-1)
	}
	return temp
}
