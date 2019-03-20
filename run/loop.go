package run

import "fmt"

type Loop struct {
	Begin   Node
	End     Node
	Step    Node
	Var     string
	IsNew   bool
	Block   []Node
	IsWhile bool
}

func (p *Module) Loop() Node {
	var loop Loop
	p.BeginScope()
	if p.check(IDENTIFIER) {
		if p.test(EQUAL) {
			loop.Var = "var_" + p.advance().Lexeme
			p.rollBack()
			loop.Begin = p.Assign()
		} else {
			loop.Begin = p.assignment()
			if p.typeOf(loop.Begin).Kind == BOOL {
				loop.IsWhile = true
				goto end
			}
		}
	}
	if loop.Var == "" {
		loop.Var = "loop_" + fmt.Sprint(p.Tokens[p.CurToken].Line)
		// if p.check(LEFT_BRACE) {
		// 	goto end
		// }
		if !p.check(RANGE) {
			loop.Begin = p.assignment()
		}
	}
	if p.match(RANGE) {
		loop.End = p.assignment()
	}
	if p.match(COMMA) {
		loop.Step = p.assignment()
	}
end:
	p.consume(LEFT_BRACE, "Expecting begin of loop block")
	p.consume(EOL, "Expecting end of line")
	p.InsideLoop = true
	loop.Block = p.block()
	p.InsideLoop = false
	p.EndScope()
	return Node{Type: LOOP, Value: loop}
}

func (t *Transpiler) Loop(node Node) {
	// beginScope()
	loop := node.Value.(Loop)
	if loop.IsWhile {
		t.file.WriteString("while(")
		t.Transpile(loop.Begin)
	} else {
		t.file.WriteString("for(")
		if loop.Begin.Type > 0 {
			if loop.Var[0:5] == "loop_" {
				t.file.WriteString("number " + loop.Var + "=")
				t.Transpile(loop.Begin)
			} else {
				t.Transpile(loop.Begin)
			}
		} else {
			t.file.WriteString("number " + loop.Var + "=0")
		}
		t.file.WriteString(";" + loop.Var + "<")
		t.Transpile(loop.End)
		t.file.WriteString(";" + loop.Var)
		if loop.Step.Type == UNARY {
			t.file.WriteString("-=")
			t.Transpile(loop.Step.Value.(Node))
		} else {
			if loop.Step.Type == 0 {
				t.file.WriteString("++")
			} else {
				t.file.WriteString("+=")
				t.Transpile(loop.Step)
			}
		}
	}
	t.file.WriteString(") {\n")
	for _, n := range loop.Block {
		t.Transpile(n)
	}
	t.file.WriteString("}\n")
	// endScope()
}
