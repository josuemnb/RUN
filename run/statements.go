package run

func (p *Module) block() []Node {
	stmt := make([]Node, 0)
	for !p.match(RIGHT_BRACE) {
		stmt = append(stmt, p.parse())
	}
	if !p.isAtEnd() && !p.match(EOL) {
		p.error("Expecting end of line", 0)
	}
	return stmt
}
