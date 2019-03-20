package run

import (
	"fmt"
	"os"
)

func (p *Module) Main() Node {
	var f Function
	f.Name = "main"
	p.BeginScope()
	f.Params = make([]Param, 0)
	if p.match(LEFT_PAREN) {
		for !p.match(RIGHT_PAREN) {
			var param Param
			n := p.advance()
			param.Name = n.Lexeme
			if p.match(DECLARE) {
				k := p.advance()
				param.Type = *p.getTypeByName(k.Lexeme)
				if param.Type.Name != "" {
					p.Scopes[p.CurScope][param.Name] = Variable{Name: param.Name, Type: param.Type}
				} else {
					p.error("Expecting a variable type", 0)
				}
			} else {
				param.Type = *p.getTypeByKind(STRING)
				p.Scopes[p.CurScope][param.Name] = Variable{Name: param.Name, Type: param.Type}
			}
			f.Params = append(f.Params, param)
			p.consume(COMMA, "Expecting ) or ,")
		}
	}
	p.consume(LEFT_BRACE, "Expeting begin of block")
	p.consume(EOL, "Expeting end of line")
	p.ActualFunction = f
	f.Body = p.block()
	p.EndScope()
	p.ActualFunction = Function{}
	return Node{Type: MAIN, Value: f}
}

func (t *Transpiler) Main(node Node) {
	t.file.WriteString("int main(")
	m := node.Value.(Function)
	l := len(m.Params)
	if l > 0 {
		t.file.WriteString("int argc, char *argv[]) {\nif((argc-1)!=" + fmt.Sprint(l) + "){\nputs(\"Error: Number of args invalid\");\nreturn -1;}\n")
		for i, p := range m.Params {
			t.file.WriteString(p.Type.Name + " var_" + p.Name)
			switch p.Type.Name {
			case "string":
				t.file.WriteString("=argv[" + fmt.Sprint(i+1) + "];\n")
			case "number":
				t.file.WriteString("=atoi(argv[" + fmt.Sprint(i+1) + "]);\n")
			case "real":
				t.file.WriteString("=atof(argv[" + fmt.Sprint(i+1) + "]);\n")
			case "bool":
				t.file.WriteString("=atoi(argv[" + fmt.Sprint(i+1) + "])>0;\n")
			default:
				println("Error: Only allowed in Main params basic kinds")
				os.Exit(-1)
			}
		}
	} else {
		t.file.WriteString(") {\n")
	}
	// p.BeginScope()
	for _, n := range m.Body {
		t.Transpile(n)
	}
	// p.EndScope()
	t.file.WriteString("return 0;\n}\n")
}

func (t *Transpiler) Finish() {
	file, err := os.OpenFile("run/libc/types.h", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		println(err.Error())
		return
	}
	for b, t := range t.Program.Types {
		if t.Kind > STRING {
			file.WriteString("#define class_" + b + "_H\n")
		}
	}
	file.Close()
	Collections.Close()
}
