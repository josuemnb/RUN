package run

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Map struct {
	Name     string
	Types    []Type
	Explicit string
}

type Array struct {
	Name     string
	Type     Type
	Explicit string
}

type Collection struct {
	Name     string
	Size     int
	Explicit string
}

var (
	Collections *os.File
)

func init() {
	var err error
	Collections, err = os.OpenFile("run/libc/collections.h", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	Collections.WriteString("#pragma once\n\n")
}

func (p *Module) Array(size int) Type {
	t := p.advance()
	cls := p.getTypeByName(t.Lexeme)
	if cls.Name == "" {
		p.error("Unkown type", 0)
	}
	if cls.Class.Name != "" {
		cls.Name = "class_" + cls.Name
	}
	value := p.getTypeByKind(cls.Kind)
	explicit := "array_" + cls.Name
	// var typ Type
	var ok bool
	typ, ok := collections[explicit]
	if !ok {
		path := "run/libc/arrays.rh"
		read, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		newContents := strings.Replace(string(read), "VALUE", cls.Real, -1)
		err = ioutil.WriteFile("run/libc/temp/"+explicit+".h", []byte(newContents), 0)
		if err != nil {
			panic(err)
		}
		Collections.WriteString("#include \"temp/" + explicit + ".h\"\n")
		typ.Name = explicit
		typ.Kind = typeIdx
		typ.Collection = ARRAY | value.Kind
		typeIdx++

		methods := make(map[string]Function)
		methods["get"] = Function{Name: "get", Params: []Param{Param{Type: *p.getTypeByKind(NUMBER)}}, Return: *value}
		methods["add"] = Function{Name: "add", Params: []Param{Param{Type: *value}}}
		methods["size"] = Function{Name: "size", Return: *p.getTypeByKind(NUMBER)}
		methods["has"] = Function{Name: "has", Params: []Param{Param{Type: *value}}, Return: *p.getTypeByKind(BOOL)}
		methods["clear"] = Function{Name: "clear"}
		typ.Class = Class{Name: typ.Name, Methods: methods}
		p.addType(typ)
		collections[explicit] = typ
	}
	return *typ
}

func (p *Module) List() Type {
	t := p.advance()
	cls := p.getTypeByName(t.Lexeme)
	if cls.Name == "" {
		p.error("Unkown type", 0)
	}
	if cls.Class.Name != "" {
		cls.Name = "class_" + cls.Name
	}
	value := p.getTypeByKind(cls.Kind)
	explicit := "list_" + cls.Name
	typ, ok := collections[explicit]
	if !ok {
		typ = new(Type)
		path := "run/libc/lists.rh"
		read, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		newContents := strings.Replace(string(read), "VALUE", cls.Real, -1)
		err = ioutil.WriteFile("run/libc/temp/"+explicit+".h", []byte(newContents), 0)
		if err != nil {
			panic(err)
		}
		Collections.WriteString("#include \"temp/" + explicit + ".h\"\n")
		typ.Name = explicit
		typ.Kind = typeIdx
		typ.Collection = LIST | value.Kind
		typeIdx++

		methods := make(map[string]Function)
		methods["get"] = Function{Name: "get", Params: []Param{Param{Type: *p.getTypeByKind(NUMBER)}}, Return: *value}
		methods["add"] = Function{Name: "add", Params: []Param{Param{Type: *value}}}
		methods["size"] = Function{Name: "size", Return: *p.getTypeByKind(NUMBER)}
		methods["has"] = Function{Name: "has", Params: []Param{Param{Type: *value}}, Return: *p.getTypeByKind(BOOL)}
		methods["clear"] = Function{Name: "clear"}
		typ.Class = Class{Name: typ.Name, Methods: methods}
		p.addType(typ)
		collections[explicit] = typ
	}
	return *typ
	// scopes[curScope][name] = Variable{Name: name, Type: typ, Array: -1}
	// return Node{Type: LIST, Value: Collection{Name: name, Size: -1, Explicit: explicit}}
}

func (p *Module) Map() Type {
	p.consume(LESS, "Expecting < for map")
	typs := make([]string, 0)
	explicit := "map"
	t := p.advance()
	cls := p.getTypeByName(t.Lexeme)
	if cls.Name == "" {
		p.error("Unkown type", 0)
	}
	// if cls.Class.Name != "" {
	// 	cls.Name = "class_" + cls.Name
	// }
	key := p.getTypeByKind(cls.Kind)
	typs = append(typs, cls.Real)
	explicit += "_" + cls.Real
	p.consume(COMMA, "Expecting , or >")
	t = p.advance()
	cls = p.getTypeByName(t.Lexeme)
	if cls.Name == "" {
		p.error("Unkown type", 0)
	}
	// if cls.Kind >= STRING {
	// 	cls.Name = "class_" + cls.Name
	// }
	value := p.getTypeByKind(cls.Kind)
	typs = append(typs, cls.Real)
	explicit += "_" + cls.Real
	p.consume(GREATER, "Expecting > for map")
	typ := new(Type)
	var ok bool
	typ, ok = collections[explicit]
	if !ok {
		typ = new(Type)
		path := "run/libc/maps.rh"
		read, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		newContents := strings.Replace(string(read), "KEY", typs[0], -1)
		if key.Kind <= STRING {
			newContents = strings.Replace(newContents, "//keyIS_COMP", "", -1)
		} else {
			newContents = strings.Replace(newContents, "//keyNOT_COMP", "", -1)
		}
		newContents = strings.Replace(newContents, "VALUE", typs[1], -1)
		if value.Kind <= STRING {
			newContents = strings.Replace(newContents, "//valueIS_COMP", "", -1)
		} else {
			newContents = strings.Replace(newContents, "//valueNOT_COMP", "", -1)
		}
		err = ioutil.WriteFile("run/libc/temp/"+explicit+".h", []byte(newContents), 0)
		if err != nil {
			panic(err)
		}
		Collections.WriteString("#include \"temp/" + explicit + ".h\"\n")
		typ.Name = explicit
		typ.Collection = MAP

		methods := make(map[string]Function)
		methods["get"] = Function{Name: "get", Params: []Param{Param{Type: *key}}, Return: *value}
		methods["put"] = Function{Name: "put", Params: []Param{Param{Type: *key}, Param{Type: *value}}, Return: *value}
		methods["size"] = Function{Name: "size", Return: *p.getTypeByKind(NUMBER)}
		methods["hasKey"] = Function{Name: "hasKey", Params: []Param{Param{Type: *key}}, Return: *p.getTypeByKind(BOOL)}
		methods["hasValue"] = Function{Name: "hasValue", Params: []Param{Param{Type: *value}}, Return: *p.getTypeByKind(BOOL)}
		methods["keys"] = Function{Name: "keys", Return: *key}
		methods["clear"] = Function{Name: "clear"}
		typ.Class = Class{Name: typ.Name, Methods: methods}
		p.addType(typ)
		collections[explicit] = typ
	}
	return *typ
	// scopes[curScope][name] = Variable{Name: name, Type: typ, Array: -1}
	// return Node{Type: MAP, Value: Collection{Name: name, Explicit: explicit}}
}

func (t *Transpiler) Collection(node Node) {
	m := node.Value.(Collection)
	t.file.WriteString(m.Explicit + " var_" + m.Name)
	if m.Size > 0 {
		t.file.WriteString("(" + fmt.Sprint(m.Size) + ")")
	}
	t.file.WriteString(";\n")
}
