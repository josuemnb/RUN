package run

type Unary struct {
	Right Node
	Op    string
}

func (p *Module) unary() Node {
	if p.match(BANG, MINUS) {
		o := p.previous()
		right := p.unary()
		return Node{Type: UNARY, Value: Unary{right, o.Lexeme}}
	}
	return p.Call(nil)
}

func (t *Transpiler) Unary(node Node) {
	u := node.Value.(Unary)
	t.file.WriteString(u.Op)
	t.Transpile(u.Right)
}
