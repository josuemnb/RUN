package run

import (
	"log"
	"os"
)

type Module struct {
	HasMain           bool
	Name              string
	Address           string
	File              *os.File
	Modules           map[string]*Module
	Functions         map[string]*Function
	Variables         map[string]*Variable
	Classes           map[string]*Class
	Collections       map[string]*Collection
	Parent            interface{}
	Type              int
	Types             map[string]*Type
	Main              *Function
	CurrentFunction   *Function
	CurrentClass      *Class
	CurrentCollection *Collection
}

func NewModule() *Module {
	m := new(Module)
	m.Variables = make(map[string]*Variable)
	m.Functions = make(map[string]*Function)
	m.Modules = make(map[string]*Module)
	m.Types = make(map[string]*Type)
	m.Classes = make(map[string]*Class)
	m.Collections = make(map[string]*Collection)
	m.registerTypes()
	return m
}

func (m *Module) registerTypes() {
	V := new(Type)
	V.Name = "void"
	V.Real = "void"
	V.Kind = VOID
	m.Types[V.Name] = V

	N := new(Type)
	N.Name = "number"
	N.Real = "number"
	N.Kind = NUMBER
	m.Types[N.Name] = N

	b := new(Type)
	b.Name = "byte"
	b.Real = "byte"
	b.Kind = BYTE
	m.Types[b.Name] = b

	S := new(Type)
	S.Name = "string"
	S.Real = "string"
	S.Kind = STRING
	m.Types[S.Name] = S

	R := new(Type)
	R.Name = "real"
	R.Real = "real"
	R.Kind = REAL
	m.Types[R.Name] = R

	B := new(Type)
	B.Name = "bool"
	B.Real = "bool"
	B.Kind = BOOL
	m.Types[B.Name] = B
}

func (m *Module) getTypeByIndex(i int) *Type {
	for _, t := range m.Types {
		if t.Kind == i {
			return t
		}
	}
	return nil
}

func (m *Module) getTypeByName(n string) *Type {
	if len(n) < 2 && m.insideCollection() == false {
		log.Fatal("Error: Type name  '" + n + "' size must be bigger than 1")
	}
	if t, ok := m.Types[n]; ok {
		return t
	}
	for mod, t := range m.Modules {
		if tp := t.getTypeByName(n); tp != nil {
			tp.Name = mod + "_" + tp.Name
			return tp
		}
	}
	return nil
}

func (m *Module) insideFunction() bool {
	return m.CurrentFunction != nil
}

func (m *Module) insideClass() bool {
	return m.CurrentClass != nil
}

func (m *Module) insideCollection() bool {
	return m.CurrentCollection != nil
}

func (m *Module) getVariable(s string) *Variable {
	if m.CurrentFunction != nil {
		if v, ok := m.CurrentFunction.Variables[s]; ok {
			return v
		}
	}
	if v, ok := m.Variables[s]; ok {
		return v
	}
	return nil
}

func (m *Module) isModule(s string) bool {
	_, ok := m.Modules[s]
	return ok
}

func (m *Module) getName(s string) (int, interface{}, bool) {
	if m.CurrentFunction != nil && m.CurrentFunction.HasVariable(s) {
		return VARIABLE, m.CurrentFunction.Variables[s], true
	} else if v, ok := m.Variables[s]; ok {
		return VARIABLE, v, ok
	} else if m.isFunction(s) {
		return FUNCTION, m.Functions[s], true
	} else if m.isModule(s) {
		return MODULE, m.Modules[s], true
	} else if m.isClass(s) {
		return CLASS, m.Classes[s], true
	} else if m.insideClass() {
		if v, ok := m.CurrentClass.Fields[s]; ok {
			return VARIABLE, v, ok
		}
	}
	return 0, nil, false
}

func (m *Module) isVariable(s string) bool {
	if m.CurrentFunction != nil && m.CurrentFunction.HasVariable(s) {
		return true
	}
	if _, ok := m.Variables[s]; ok {
		return true
	}
	return false
}

func (m *Module) isFunction(s string) bool {
	if m.insideClass() {
		if _, ok := m.CurrentClass.Functions[s]; ok {
			return ok
		}
	}
	_, ok := m.Functions[s]
	return ok
}

func (m *Module) getFunction(s string) (*Function, bool) {
	if m.insideClass() {
		if f, ok := m.CurrentClass.Functions[s]; ok {
			return f, ok
		}
	}
	f, ok := m.Functions[s]
	return f, ok
}

func (m *Module) isClass(s string) bool {
	_, ok := m.Classes[s]
	return ok
}
