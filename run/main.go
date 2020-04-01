package run

import (
	"fmt"
	"log"
)

func (m *Module) parseMain(node *Node) {
	if m.Main != nil {
		println("Error: Module already has a main function at " + fmt.Sprint(node.Token.Line))
	}
	m.Main = NewFunction()
	m.Main.Parent = m
	m.Main.Name = "main"
	m.Main.Real = "main"
	m.Main.Return = NUMBER
	m.CurrentFunction = m.Main

	node.Code = "int main("
	if len(node.Children) == 0 {
		log.Fatal("Error: Expecting body or params for main at " + fmt.Sprint(node.Token.Line))
	}
	node.Parsed = true
	n := node.Children[0]
	l := len(n.Children)
	if n.Type == LEFT_PAREN {
		if node.Children[1].Type != RIGHT_PAREN {
			log.Fatal("Error: End of declaration not correct at " + fmt.Sprint(node.Token.Line))
		}
		if l > 0 {
			n.Parsed = true
			m.parseMainParams(n)
			n.Code = "int argc, char *argv[]) {\n\t" +
				"if((argc-" + fmt.Sprint(len(m.Main.Params)) + ")!=1){\n\t" +
				"puts(\"Error: Number of args invalid\");\n\t" +
				"return -1;\n}\n"
		} else {
			n.Code = ") {\n"
		}
		n = node.Children[2]
	} else {
		node.Code += ") {\n"
	}
	l = len(node.Children)
	if n.Type != LEFT_BRACE {
		log.Fatal("Error: Expecting begin of body at " + fmt.Sprint(node.Token.Line))
	}
	m.parseBody(n)
	if l <= n.Index+1 || node.Children[n.Index+1].Type != RIGHT_BRACE {
		log.Fatal("Error: Expecting end of body at " + fmt.Sprint(node.Token.Line))
	}
	node.Children[n.Index+1].Code = "\nreturn 0;\n}\n"
	m.CurrentFunction = nil
}

func (m *Module) parseMainParams(params *Node) {
	l := len(params.Children)
	index := 0
	counter := 1
	for {
		index = m.parseMainParamsNode(params, index, counter)
		counter++
		if index >= l {
			break
		}
		if params.Children[index].Type != COMMA {
			log.Fatal("Error: Expecting ',' at " + fmt.Sprint(params.Token.Line))
		}
		params.Children[index].Parsed = true
		index++
	}
}

func (m *Module) parseMainParamsNode(params *Node, index int, number int) int {
	l := len(params.Children)
	n := params.Children[index]
	n.Parsed = true
	if n.Type != IDENTIFIER {
		log.Fatal("Error: Expecting a Identifier at " + fmt.Sprint(params.Token.Line))
	}
	if m.Main.HasParam(n.Token.Value) {
		log.Fatal("Error: Name already in use at " + fmt.Sprint(params.Token.Line))
	}
	if l == index+1 {
		n.Code = "Run_string " + n.Token.Value + " = argv[" + fmt.Sprint(number) + "];\n"
		m.Main.AddParam(&Variable{Name: n.Token.Value, Type: m.getTypeByName("string")})
		return index + 1
	}
	if l < index+3 {
		log.Fatal("Error: Expecting a type at " + fmt.Sprint(params.Token.Line))
	}
	n1 := params.Children[index+1]
	if n1.Type == COMMA {
		n.Code = "Run_string " + n.Token.Value + " = argv[" + fmt.Sprint(number) + "];\n"
		m.Main.AddParam(&Variable{Name: n.Token.Value, Type: m.getTypeByName("string")})
		return index + 1
	}
	n1.Parsed = true
	if n1.Type != DECLARE {
		log.Fatal("Error: Expecting a ':'' at " + fmt.Sprint(params.Token.Line))
	}
	n1 = params.Children[index+2]
	n1.Parsed = true
	if n.Type != IDENTIFIER {
		log.Fatal("Error: Expecting a Type at " + fmt.Sprint(params.Token.Line))
	}
	var f string
	var b string
	if n1.Token.Value == "number" {
		f = "func_number_string("
		b = ")"
	}
	n.Code = stringToType(n1.Token.Value) + " " + n.Token.Value + " = " + f + "argv[" + fmt.Sprint(number) + "]" + b + ";\n"
	m.Main.AddParam(&Variable{Name: n.Token.Value, Type: m.getTypeByName(stringToType(n1.Token.Value))})
	return index + 3
}
