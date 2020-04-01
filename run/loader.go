package run

import (
	"io/ioutil"
	"log"
	"os"
)

func LoadModule(prog string) *Module {
	return load(prog, false)
}

func LoadProgram(prog string) *Module {
	return load(prog, true)
}

func load(prog string, isMain bool) *Module {
	module := NewModule()
	module.Name = prog

	if _, err := os.Stat(prog); os.IsNotExist(err) {
		log.Fatal("Error: File not found '" + module.Name + "'")
	}
	b, err := ioutil.ReadFile(prog)
	if err != nil {
		log.Fatal(err.Error())
	}
	if isMain {
		module.Type = Types.MAIN
	} else {
		module.Type = Types.MODULE
		l := (7 + len(module.Name))
		if string(b[0:l]) != (Values.MODULE + module.Name) {
			log.Fatal("Error: Module incorret")
		}
	}

	module.File, err = os.OpenFile(module.Name+".cc", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	if isMain {
		module.Name = "Root"
	}
	module.File.WriteString("#include \"run/libc/run.h\"\n\n")
	scanner := NewScanner(b)
	scanner.scan()
	if isMain && scanner.HasMain == false {
		log.Println("Warning: Main not found!")
		// os.Exit(-2)
	}
	scanner.printTree()
	module.Parse(scanner.Root)
	PrintCode(module.File, scanner.Root)
	return module
}

func PrintCode(file *os.File, node *Node) {
	if node.NotPrint {
		return
	}
	if node.Code != "" {
		file.WriteString(node.Code)
	}
	for _, v := range node.Children {
		PrintCode(file, v)
	}
}
