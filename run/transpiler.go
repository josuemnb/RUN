package run

import (
	"os"
)

type Transpiler struct {
	file        *os.File
	class       Class
	insideClass bool
	Program     *Module
	Printing    bool
}

func NewTranspiler(Program *Module) *Transpiler {
	t := new(Transpiler)
	t.file = Program.File
	t.Program = Program
	return t
}

func (t *Transpiler) Transpile(node Node) {
	switch node.Type {
	case NEWLINE:
		t.file.WriteString("\"\\n\"")
	case NULL, NIL:
		t.file.WriteString("NULL")
	case MAP, ARRAY, LIST, STACK:
		t.Collection(node)
	case BRACKETS:
		t.Bracket(node)
	case CPP:
		t.Cpp(node)
	case DOT:
		t.Dot(node)
	case ASSIGN:
		t.Assign(node)
	case LITERAL:
		t.Literal(node)
	case BINARY:
		t.Binary(node)
	case UNARY:
		t.Unary(node)
	case GROUPING:
		t.Grouping(node)
	case IF, ELSEIF:
		t.If(node)
	case IDENTIFIER:
		t.Identifier(node)
	case FUNCTION:
		t.Function(node)
	case RETURN:
		t.Return(node)
	case THIS:
		t.This(node)
	case EOL:
		t.file.WriteString(";\n")
	case DECLARE:
		t.Declare(node)
	case PRINT:
		t.Print(node)
	case CALL:
		t.Call(node)
	case MAIN:
		t.Main(node)
	case LOOP:
		t.Loop(node)
	case BREAK:
		// 	// println("break;")
		t.file.WriteString("break;\n")
	case CLASS:
		t.Class(node)
	case BETWEEN:
		t.Between(node)
	}
}
