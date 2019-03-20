package main

import (
	"./run"
)

func main() {
	prog := run.NewProgram("test.run")
	prog.Compile("test.cpp")
	prog.Finish()
}
