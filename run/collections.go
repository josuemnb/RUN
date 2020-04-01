package run

import (
	"fmt"
	"log"
)

type Collection struct {
	*Class
	Values []string
	Node   *Node
}

func NewCollection() *Collection {
	c := new(Collection)
	c.Values = make([]string, 0)
	c.Class = NewClass()
	return c
}

func (m *Module) parseCollection(node *Node) {
	l := len(node.Children)
	if l == 1 || node.Children[1].Type != RIGHT_BRACKET {
		log.Fatal("Error: Badformed line at " + fmt.Sprint(node.Token.Line))
	}
	if l > 3 {
		if node.Children[2].Type == LEFT_BRACE && node.Children[3].Type == RIGHT_BRACE {
			m.parseCollectionDeclare(node)
		} else if node.Children[2].Type == EQUAL {
			node.Code = node.Token.Value
			node.Children[2].Code = " = "
			node.Children[0].Code = "["
			if len(node.Children[0].Children) == 0 || node.Children[0].Children[0].Type != NUMBER {
				log.Fatal("Error: Expecting an index at " + fmt.Sprint(node.Token.Line))
			}
			node.Children[0].Children[0].Code = node.Children[0].Children[0].Token.Value
			node.Children[1].Code = "]"
			v := m.getVariable(node.Token.Value)
			index, typ := m.parseExpression(node, 3)
			if l <= index || node.Children[index].Type != NEWLINE {
				log.Fatal("Error: Expecting end of line at " + fmt.Sprint(node.Token.Line))
			}
			node.Children[index].Code = ";\n"
			if v.Type.Kind != typ {
				log.Fatal("Error: Badformed line at " + fmt.Sprint(node.Token.Line))
			}
		} else if node.Children[2].Type == DOT {

		}
	}
}

func (m *Module) parseCollectionDeclare(node *Node) {
	coll := NewCollection()
	l := len(node.Children[0].Children)
	index := 0
	coll.Name = m.Name + "_" + node.Token.Value
	coll.Node = node
	// node.NotPrint = true
	node.Code = Values.CLASS + " " + coll.Name + "_"
	count := 1_000_000
	for {
		n := node.Children[0].Children[index]
		if n.Type != IDENTIFIER {
			log.Fatal("Error: Expecting ValuedType at " + fmt.Sprint(node.Token.Line))
		}
		if len(n.Token.Value) != 1 {
			log.Fatal("Error: Expected ValuedType Letter, not a string at " + fmt.Sprint(node.Token.Line))
		}
		n.Code = n.Token.Value
		coll.Name += "_" + n.Token.Value
		coll.Values = append(coll.Values, n.Token.Value)
		m.Types[n.Token.Value] = &Type{Kind: count, Name: n.Token.Value, Collection: coll, Real: n.Token.Value}
		index++
		if index >= l {
			break
		}
		if node.Children[0].Children[index].Type != COMMA {
			log.Fatal("Error: Expecting comma at " + fmt.Sprint(node.Token.Line))
		}
		node.Children[0].Children[index].Code = "_"
		count++
		index++
	}
	coll.Real = coll.Name
	m.Collections[coll.Name] = coll
	m.CurrentCollection = coll
	m.CurrentClass = coll.Class
	if len(node.Children[2].Children) > 0 {
		m.parseBody(node.Children[2])
	}
	m.CurrentCollection = nil
	m.CurrentClass = nil
	node.Children[2].Code = " {\n"
	node.Children[3].Code = "\n};\n\n"
}

func (c *Collection) isValuedType(s string) bool {
	for _, k := range c.Values {
		if k == s {
			return true
		}
	}
	return false
}
