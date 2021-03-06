package run

type Cpp struct {
	Body string
}

func (p *Module) Cpp() Node {
	// var buff bytes.Buffer
	// p.consume(MESSAGE, "Expecting #")
	// braces := 0
	// for {
	// 	if p.check(LEFT_BRACE) {
	// 		braces++
	// 	} else if p.check(RIGHT_BRACE) {
	// 		if braces <= 0 {
	// 			break
	// 		}
	// 		braces--
	// 	}
	// 	tok := p.Tokens[p.CurToken]
	// 	if tok.Type == STRING {
	// 		buff.WriteString("\"")
	// 	}
	// 	buff.WriteString(tok.Lexeme)
	// 	if tok.Type == STRING {
	// 		buff.WriteString("\"")
	// 	} else {
	// 		buff.WriteString(" ")
	// 	}
	// 	p.CurToken++
	// }
	// p.consume(CARDINAL, "Expecting #")
	// println(buff.String())
	s := p.consume(MESSAGE, "Expecting #")
	// println("CPP", s.Lexeme)
	return Node{Type: CPP, Value: Cpp{s.Lexeme}}
}

func (t *Transpiler) Cpp(node Node) {
	t.file.WriteString(node.Value.(Cpp).Body)
}
