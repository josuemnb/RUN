package run

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func (m *Module) parsePrint(node *Node, newLine bool) {
	var buff strings.Builder
	// var back strings.Builder
	node.Parsed = true
	l := len(node.Children)
	if l == 0 {
		log.Fatal("Error: Expecting '(' at " + fmt.Sprint(node.Token.Line))
		os.Exit(-2)
	}
	params := node.Children[0]
	if params.Type != LEFT_PAREN {
		log.Fatal("Error: Expecting '(' at " + fmt.Sprint(node.Token.Line))
		os.Exit(-2)
	}
	idx := l - 1
	if node.Children[idx].Type == NEWLINE {
		idx--
	}
	if node.Children[idx].Type != RIGHT_PAREN {
		log.Fatal("Error: Expecting ')' at " + fmt.Sprint(node.Token.Line))
		os.Exit(-2)
	}
	node.Children[idx].Code = ");\n"
	params.Parsed = true
	l = len(params.Children)
	index := 0
	var s string
	buff.WriteString("printf(\"")
	for {
		s, index = m.parsePrintParam(params, index)
		buff.WriteString(s)
		if index >= l {
			break
		}
	}
	if newLine {
		buff.WriteString("\\n")
	}
	params.Code = buff.String() + "\", " // + back.String() + ");\n"
}

func stringToPrint(v string) (s string) {
	switch v {
	case "Run_number":
		s = "%lld"
	case "Run_bool":
		s = "%s"
	case "Run_real":
		s = "%f"
	case "Run_string":
		s = "%s"
	case "Run_byte":
		s = "%c"
	}
	return
}

func intToPrint(t int) string {
	switch t {
	case NUMBER:
		return "%lld"
	case REAL:
		return "%f"
	case BYTE:
		return "%c"
	case STRING, QUOTE, BOOL:
		return "%s"
	}
	return "%x"
}

func (m *Module) parsePrintParam(node *Node, index int) (s string, i int) {
	node.Parsed = true
	if node.Children[index].Type == COMMA {
		s = " "
		node.Code = ", "
		i = index + 1
		return
	}
	var t int
	i, t = m.parseExpression(node, index)
	if index > 0 {
		node.Children[index].Code = ", " + node.Children[index].Code
	}
	s += intToPrint(t)
	return
}
