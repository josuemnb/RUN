package run

import (
	"fmt"
	"log"
	"os"
)

func (m *Module) parseExpression(node *Node, index int) (int, int) {
	temp, t := m.parseOr(node, index)
	// if len(node.Children) > temp {
	// 	n := node.Children[temp]
	// 	if n.Type == EQUAL {
	// 		n.Code = " = "
	// 		n.Parsed = true
	// 		temp++
	// 		return m.parseExpression(node, temp)
	// 	}
	// }
	return temp, t
}

func (m *Module) parseOr(node *Node, index int) (int, int) {
	temp, t := m.parseAnd(node, index)
	if len(node.Children) > temp {
		n := node.Children[temp]
		if n.Type == OR {
			temp++
			n.Parsed = true
			n.Code = " || "
			return m.parseOr(node, temp)
		}
	}
	return temp, t
}

func (m *Module) parseAnd(node *Node, index int) (int, int) {
	temp, t := m.parseEqual(node, index)
	if len(node.Children) > temp {
		n := node.Children[temp]
		if n.Type == AND {
			temp++
			n.Code = " && "
			return m.parseAnd(node, temp)
		}
	}
	return temp, t
}

func (m *Module) parseEqual(node *Node, index int) (int, int) {
	temp, t := m.parseComparation(node, index)
	if len(node.Children) > temp {
		n := node.Children[temp]
		if n.Type == BANG_EQUAL || n.Type == EQUAL_EQUAL {
			n.Code = n.Token.Value
			temp++
			n.Parsed = true
			return m.parseEqual(node, temp)
		}
	}
	return temp, t
}

func (m *Module) parseComparation(node *Node, index int) (int, int) {
	temp, t := m.parseAddition(node, index)
	if len(node.Children) > temp {
		n := node.Children[temp]
		if n.Type == GREATER || n.Type == GREATER_EQUAL || n.Type == LESS || n.Type == LESS_EQUAL {
			temp++
			n.Parsed = true
			n.Code = n.Token.Value
			return m.parseComparation(node, temp)
		}
	}
	return temp, t
}

func (m *Module) parseAddition(node *Node, index int) (int, int) {
	temp, t := m.parseMultiplication(node, index)
	if len(node.Children) > temp {
		n := node.Children[temp]
		if n.Type == PLUS || n.Type == MINUS {
			temp++
			n.Parsed = true
			n.Code = n.Token.Value
			return m.parseAddition(node, temp)
		}
	}
	return temp, t
}

func (m *Module) parseMultiplication(node *Node, index int) (int, int) {
	temp, t := m.parseIncrement(node, index)
	if len(node.Children) > temp {
		n := node.Children[temp]
		if n.Type == MUL || n.Type == DIV {
			temp++
			n.Parsed = true
			n.Code = n.Token.Value
			return m.parseMultiplication(node, temp)
		}
	}
	return temp, t
}

func (m *Module) parseIncrement(node *Node, index int) (int, int) {
	temp, t := m.parseUnary(node, index)
	if len(node.Children) > temp {
		n := node.Children[temp]
		if n.Type == INCREMENT || n.Type == DECREMENT {
			temp++
			n.Parsed = true
			n.Code = n.Token.Value
			return m.parseIncrement(node, temp)
		}
	}
	return temp, t
}

func (m *Module) parseUnary(node *Node, index int) (int, int) {
	// parseEqual(node, index)
	if len(node.Children) > index {
		n := node.Children[index]
		if n.Type == MINUS || n.Type == PLUS {
			index++
			n.Parsed = true
			n.Code = n.Token.Value
			return 0, index
		}
	}
	return m.parseCall(node, index)
}

func (m *Module) parseCall(node *Node, index int) (int, int) {
	temp, t := m.parsePrimary(node, index)
	return temp, t
}

func (m *Module) parsePrimary(node *Node, index int) (int, int) {
	temp := index
	typ := 0
	if len(node.Children) > index {
		n := node.Children[index]
		n.Parsed = true
		if n.Type == BOOL {
			n.Code = n.Token.Value
			typ = n.Type
			temp++
			goto end
			// return index + 1, n.Type
		} else if n.Type == NIL {

		} else if n.Type == NUMBER || n.Type == QUOTE || n.Type == REAL || n.Type == BOOL {
			n.Code = n.Token.Value
			if n.Type == QUOTE {
				// return index + 1, m.getTypeByIndex(STRING).Kind
				temp++
				typ = m.getTypeByIndex(STRING).Kind
				goto end
			}
			// return index + 1, n.Type
			temp++
			typ = n.Type
			goto end
		} else if n.Type == IDENTIFIER {
			if v, ok := m.CurrentFunction.Variables[n.Token.Value]; ok {
				n.Code = n.Token.Value
				// return index + 1, v.Type.Kind
				temp++
				typ = v.Type.Kind
				if v.isArray() {
					m.checkArray(v, n)
				}
				if typ >= 1000 {
					m.checkClass(v, n)
				}
			}
		} else if n.Type == LEFT_PAREN {
			n.Code = "("
			_, t := m.parseExpression(n, 0)
			// n.Children[i-1].Code += ")"
			if node.Children[index+1].Type != RIGHT_PAREN {
				log.Fatal("Error: Expecting ')' at " + fmt.Sprint(node.Token.Line))
				os.Exit(-1)
			}
			node.Children[index+1].Code = ")"
			// return index + 2, t
			temp += 2
			typ = t
			goto end
		}
	}
end:
	return temp, typ
}

func (m *Module) checkArray(v *Variable, node *Node) {
	l := len(node.Children)
	if l > 1 && node.Children[0].Type == LEFT_BRACKET && node.Children[1].Type == RIGHT_BRACKET {
		node.Children[0].Code = "["
		node.Children[1].Code = "]"
		node.Children[0].Parsed = true
		node.Children[1].Parsed = true
		_, typ := m.parseExpression(node.Children[0], 0)
		if typ != NUMBER {
			log.Fatal("Error: Expecting numbered value at " + fmt.Sprint(node.Token.Line))
		}
	}
}

func (m *Module) checkClass(v *Variable, node *Node) {

}
