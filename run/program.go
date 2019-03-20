package run

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Module struct {
	Name           string
	File           *os.File
	Transpiler     Transpiler
	Modules        map[string]*Module
	Parent         *Module
	Functions      map[string]Function
	Types          map[string]*Type
	HasMain        bool
	Type           int
	Scopes         []map[string]Variable
	CurScope       int
	Tokens         []Token
	CurToken       int
	ActualClass    Class
	ActualFunction Function
	InsideLoop     bool
	Ignore         bool
}

func (p *Module) insideFunction() bool {
	return p.ActualFunction.Name != ""
}

func (p *Module) insideClass() bool {
	return p.ActualClass.Name != ""
}

func NewProgram(name string) *Module {
	program := new(Module)
	program.Name = name
	// program.Classes = make(map[string]Class)
	program.Types = make(map[string]*Type)
	program.Functions = make(map[string]Function)
	program.Modules = make(map[string]*Module)
	program.Scopes = make([]map[string]Variable, 0)
	program.Scopes = append(program.Scopes, make(map[string]Variable))

	V := new(Type)
	V.Name = "void"
	V.Kind = VOID
	program.addType(V)

	N := new(Type)
	N.Name = "number"
	N.Kind = NUMBER
	program.addType(N)

	program.addType(&Type{Name: "quote", Kind: QUOTE})

	S := new(Type)
	S.Name = "string"
	S.Kind = STRING
	program.addType(S)
	S.Class.Methods = make(map[string]Function)
	S.Class.Methods["substring_number_number"] = Function{Name: "substring", Real: "substring_number_number", Params: []Param{Param{Type: *program.getTypeByKind(NUMBER)}, Param{Type: *program.getTypeByKind(NUMBER)}}, Return: *program.getTypeByKind(STRING)}
	S.Class.Methods["size"] = Function{Name: "size", Real: "size", Return: *program.getTypeByKind(NUMBER)}

	R := new(Type)
	R.Name = "real"
	R.Kind = REAL
	program.addType(R)

	B := new(Type)
	B.Name = "bool"
	B.Kind = BOOL
	program.addType(B)

	typeIdx = 132
	return program
}

func (p *Module) Compile(name string) {
	n := p.Name
	if p.Type == MODULE {
		n = "run/lib/" + n + ".rmd"
	}
	if _, err := os.Stat(n); os.IsNotExist(err) {
		log.Fatal("Error: File not found '" + name + "'")
		os.Exit(-1)
	}
	b, err := ioutil.ReadFile(n)
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	p.File, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	var s *Scanner
	if p.Type == MODULE {
		l := (7 + len(p.Name))
		if string(b[0:l]) != ("module " + p.Name) {
			println("Error: Module incorret")
			os.Exit(-1)
		}
		s = NewScanner(b[l:], true)
		p.File.WriteString("#include \"../../libc/run.h\"\n")
	} else {
		s = NewScanner(b, false)
		p.File.WriteString("#include \"run/libc/run.h\"\n")
	}
	p.Tokens = s.ScanTokens()
	if s.HasMain() == false {
		println("Error: Main not found")
		os.Exit(-1)
	}
	p.Transpiler = *NewTranspiler(p)
	for _, node := range p.Parse() {
		p.Transpiler.Transpile(node)
	}
	p.File.WriteString(";\n")
	p.File.Close()
}

func (p *Module) Import() Node {
	name := p.consume(IDENTIFIER, "Expecting name")
	m := "run/lib/temp/module_" + name.Lexeme + ".h"
	p.File.WriteString("#include \"" + m + "\"\n")
	// if _, err := os.Stat(m); err == nil {
	// 	p.File.WriteString("#include \"" + m + "\"\n")
	// 	return Node{}
	// }
	prog := NewProgram(name.Lexeme)
	prog.Type = MODULE
	prog.HasMain = true
	prog.Compile(m)
	prog.Finish()
	// p.File.WriteString("#include \"" + m + "\"\n")
	p.Modules[name.Lexeme] = prog
	return Node{}
}

func (p *Module) Finish() {
	p.Transpiler.Finish()
	compile(p.Name[:len(p.Name)-4])
}

func compile(arg string) bool {
	// println("Compiling...")
	exec.Command("gcc", "-o", arg, arg+".cpp", "-O3", "-w").Output()
	os.Remove(arg + ".cpp")
	// println("Done")
	return true
}
