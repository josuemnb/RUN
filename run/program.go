package run

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Module struct {
	Name       string
	File       *os.File
	Transpiler Transpiler
	Modules    map[string]*Module
	Parent     *Module
	Functions  map[string]Function
	Types      map[string]*Type
	// Interfaces     map[string]Interface
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
	Link           string
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
	program.Type = MAIN
	program.Types = make(map[string]*Type)
	program.Functions = make(map[string]Function)
	program.Modules = make(map[string]*Module)
	// program.Interfaces = make(map[string]Interface)
	program.Scopes = make([]map[string]Variable, 0)
	program.Scopes = append(program.Scopes, make(map[string]Variable))

	V := new(Type)
	V.Name = "void"
	V.Kind = VOID
	program.addType(V)

	N := new(Type)
	N.Name = "number"
	N.Real = "RUN_"
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
	R.Real = "RUN_"
	R.Kind = REAL
	program.addType(R)

	B := new(Type)
	B.Name = "bool"
	B.Real = "RUN_"
	B.Kind = BOOL
	program.addType(B)

	// typeIdx = 132
	return program
}

func (p *Module) Compile() {
	n := p.Name
	if p.Type == MODULE {
		n = "run/lib/" + n + ".rmd"
	} else {
		n += ".run"
	}
	if _, err := os.Stat(n); os.IsNotExist(err) {
		log.Fatal("Error: File not found '" + p.Name + "'")
		os.Exit(-1)
	}
	b, err := ioutil.ReadFile(n)
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	var f string
	if p.Type == MODULE {
		f = "run/lib/temp/module_" + p.Name + ".h"
	} else {
		f = p.Name + ".cpp"
	}

	p.File, err = os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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
		s = NewScanner(b, true)
		p.File.WriteString("#pragma once\n\n#include \"../../libc/run.h\"\n")
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
	if p.Type != MODULE {
		p.Finish(true)
	}
}

func (p *Module) Module() Node {
	n := p.consume(IDENTIFIER, "Expecting name of module").Lexeme
	if n != p.Name {
		p.error("Name of module different", 0)
	}
	if !p.match(EOL) {
		p.Link = p.consume(QUOTE, "Expecting link value").Lexeme
	}
	return Node{}
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
	prog.Compile()
	// prog.Finish(false)
	// p.File.WriteString("#include \"" + m + "\"\n")
	p.Modules[name.Lexeme] = prog
	if p.match(COMMA) {
		return p.Import()
	}
	return Node{}
}

func (p *Module) Finish(comp bool) {
	p.Transpiler.Finish()
	if p.Type != MODULE && comp {
		p.compile(p.Name)
	}
	if comp {
		// os.Remove(p.Name + ".cpp")
	}
}

func (p *Module) compile(arg string) bool {
	println("Compiling...")
	link := ""
	for _, m := range p.Modules {
		if m.Link != "" {
			if link != "" {
				link += " "
			}
			link += m.Link
		}
	}
	var cmd *exec.Cmd
	if len(link) > 0 {
		cmd = exec.Command("gcc", "-o", arg, arg+".cpp", "-O3", "-s", "-w", "-std=c99", link)
	} else {
		cmd = exec.Command("gcc", "-o", arg, arg+".cpp", "-O3", "-s", "-w", "-std=c99")
	}
	err := cmd.Run()
	if err != nil {
		println("Error", err.Error())
	} else {
		//
		println("Done")
	}
	return true
}
