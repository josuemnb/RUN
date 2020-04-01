package run

import (
	"fmt"
	"log"
)

type Class struct {
	Name            string
	Real            string
	Fields          map[string]*Variable
	Functions       map[string]*Function
	Parent          interface{}
	This            []*Function
	Type            *Type
	CurrentFunction *Function
}

func NewClass() *Class {
	c := new(Class)
	c.Fields = make(map[string]*Variable)
	c.Functions = make(map[string]*Function)
	c.This = make([]*Function, 0)
	return c
}

func (m *Module) parseClassDeclare(node *Node) {
	class := NewClass()
	class.Parent = m
	class.Name = node.Token.Value
	class.Real = m.Name + "_" + class.Name
	m.Classes[node.Token.Value] = class
	class.Type = &Type{Name: class.Name, Kind: len(m.Types) + 1000, Real: class.Real, Module: m, Class: class}
	m.Types[class.Name] = class.Type
	m.CurrentClass = class
	if len(node.Children) != 2 || node.Children[1].Type != RIGHT_BRACE {
		log.Fatal("Error: Expecting end of class at " + fmt.Sprint(node.Token.Line))
	}
	node.Code = Values.CLASS + class.Real + " {\n"
	node.Parsed = true
	m.parseBody(node.Children[0])
	node.Children[1].Code = "};\n\n"
	node.Children[1].Parsed = true
	m.CurrentClass = nil
}

func (m *Module) isClassThis(s string) bool {
	for _, f := range m.CurrentClass.This {
		if f.Real == s {
			return true
		}
	}
	return false
}

func (m *Module) parseThis(node *Node) {
	if m.insideClass() == false {
		log.Fatal("Error: this only inside class " + fmt.Sprint(node.Token.Line))
	}
	node.Code = Values.PUBLIC + m.CurrentClass.Real
	l := len(node.Children)
	if m.insideFunction() {
		if l == 0 {

		} else {

		}
	} else {
		if l < 4 {
			log.Fatal("Error: 'this' must be declare as function at " + fmt.Sprint(node.Token.Line))
		}
		if node.Children[0].Type != LEFT_PAREN || node.Children[1].Type != RIGHT_PAREN || node.Children[2].Type != LEFT_BRACE || node.Children[3].Type != RIGHT_BRACE {
			log.Fatal("Error: 'this' must be declare as function at " + fmt.Sprint(node.Token.Line))
		}
		function := NewFunction()
		function.Name = node.Token.Value
		function.Real = node.Token.Value
		function.Parent = m.CurrentClass
		m.CurrentFunction = function
		m.CurrentClass.CurrentFunction = function
		node.Children[0].Code = "("
		node.Children[1].Code = ") "
		node.Children[2].Code = "{\n"
		node.Children[3].Code = "}\n"
		if len(node.Children[0].Children) > 0 {
			m.parseParams(node.Children[0])
		}
		if m.isClassThis(function.Real) {
			log.Fatal("Error: 'this' already assigned at " + fmt.Sprint(node.Token.Line))
		}
		m.CurrentClass.This = append(m.CurrentClass.This, function)
	}
	m.CurrentFunction = nil
	m.CurrentClass.CurrentFunction = nil
}
