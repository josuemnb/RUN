package run

import (
	"fmt"
	"log"
	"strings"
)

type Function struct {
	Name      string
	Address   string
	Real      string
	Parent    interface{}
	Params    []*Variable
	Variables map[string]*Variable
	Return    int
	Access    int
}

func NewFunction() *Function {
	f := new(Function)
	f.Params = make([]*Variable, 0)
	f.Variables = make(map[string]*Variable)
	f.Access = PUBLIC
	return f
}

func (f *Function) HasParam(s string) bool {
	for _, p := range f.Params {
		if s == p.Name {
			return true
		}
	}
	return false
}

func (f *Function) HasVariable(s string) bool {
	for _, p := range f.Variables {
		if s == p.Name {
			return true
		}
	}
	return false
}

func (f *Function) AddParam(v *Variable) {
	f.Params = append(f.Params, v)
	f.Variables[v.Name] = v
	f.Real += "_" + v.Type.Name
}

func (m *Module) FunctionDeclare(node *Node) {
	index := 0
	if node.Children[index+1].Type != RIGHT_PAREN {
		log.Fatal("Error: End of declaration not correct at " + fmt.Sprint(node.Token.Line))
	}
	name := node.Token.Value
	node.Parsed = true
	function := NewFunction()
	function.Return = VOID
	function.Parent = m
	function.Name = name
	m.CurrentFunction = function
	if m.insideClass() {
		if strings.HasPrefix(name, "_") {
			if strings.HasPrefix(name, "__") {
				node.Code = Values.PRIVATE
				function.Access = PRIVATE
			} else {
				node.Code = Values.PROTECTED
				function.Access = PROTECTED
			}
		} else {
			node.Code = Values.PUBLIC
			function.Access = PUBLIC
		}
		m.CurrentClass.Functions[name] = function
		m.CurrentClass.CurrentFunction = function
		function.Real = m.CurrentClass.Real + "_" + function.Name
	} else {
		m.Functions[name] = function
		function.Real = m.Name + "_" + name
	}
	if len(node.Children[index].Children) > 0 {
		m.parseParams(node.Children[index])
	}
	index += 2
	if node.Children[index].Type == DECLARE {
		index++
		if node.Children[index].Type != IDENTIFIER {
			log.Fatal("Error: Expecting type of return at " + fmt.Sprint(node.Token.Line))
		}
		typ := m.getTypeByName(node.Children[index].Token.Value)
		if typ == nil {
			log.Fatal("Error: Unknown type '" + node.Children[index].Token.Value + "' at " + fmt.Sprint(node.Token.Line))
		}
		function.Return = typ.Kind
		index++
	}
	if node.Children[index].Type != LEFT_BRACE {
		log.Fatal("Error: Expecting '{' at " + fmt.Sprint(node.Token.Line))
	}
	if node.Children[index+1].Type != RIGHT_BRACE {
		log.Fatal("Error: Expecting '}' at " + fmt.Sprint(node.Token.Line))
	}
	node.Code += m.getTypeByIndex(m.CurrentFunction.Return).Real + " " + name + "("
	node.Children[index].Code = ") {\n"
	if len(node.Children[index].Children) > 0 {
		m.parseBody(node.Children[index])
	}
	if function.Return != VOID {
		node.Children[index+1].Code = "\nterminate(\"Error: Function returning value undefined\");\n}\n"
	} else {
		node.Children[index+1].Code += "\n}\n"
	}
	m.CurrentFunction = nil
	if m.insideClass() {
		m.CurrentClass.CurrentFunction = nil
	}
}

func (m *Module) parseParams(params *Node) {
	l := len(params.Children)
	index := 0
	counter := 1
	for {
		index = m.parseParamsNode(params, index, counter)
		counter++
		if index >= l {
			break
		}
		if params.Children[index].Type != COMMA {
			log.Fatal("Error: Expecting ',' at " + fmt.Sprint(params.Token.Line))
		}
		params.Children[index].Parsed = true
		params.Children[index].Code = ", "
		index++
	}
}

func (m *Module) parseParamsNode(params *Node, index int, number int) int {
	l := len(params.Children)
	n := params.Children[index]
	n.Parsed = true
	if n.Type != IDENTIFIER {
		log.Fatal("Error: Expecting a Identifier at " + fmt.Sprint(params.Token.Line))
	}
	if m.CurrentFunction.HasParam(n.Token.Value) {
		log.Fatal("Error: Name already in use at " + fmt.Sprint(params.Token.Line))
	}
	if l == index+1 {
		n.Code = "string " + n.Token.Value
		m.CurrentFunction.AddParam(&Variable{Name: n.Token.Value, Type: m.getTypeByName("string")})
		return index + 1
	}
	if l < index+3 {
		log.Fatal("Error: Expecting a type at " + fmt.Sprint(params.Token.Line))
	}
	n1 := params.Children[index+1]
	if n1.Type == COMMA {
		n.Code = "string " + n.Token.Value
		m.CurrentFunction.AddParam(&Variable{Name: n.Token.Value, Type: m.getTypeByName("string")})
		return index + 1
	}
	n1.Parsed = true
	if n1.Type != DECLARE {
		log.Fatal("Error: Expecting a ':'' at " + fmt.Sprint(params.Token.Line))
	}
	n1 = params.Children[index+2]
	n1.Parsed = true
	if n.Type != IDENTIFIER {
		log.Fatal("Error: Expecting a Type at " + fmt.Sprint(params.Token.Line))
	}
	n.Code = stringToType(n1.Token.Value) + " " + n.Token.Value
	m.CurrentFunction.AddParam(&Variable{Name: n.Token.Value, Type: m.getTypeByName(n1.Token.Value)})
	return index + 3
}

func (m *Module) parseReturn(node *Node) {
	if m.insideFunction() == false {
		log.Fatal("Error: Return only allowed inside functions body at " + fmt.Sprint(node.Token.Line))
	}
	if m.CurrentFunction.Return == VOID {
		log.Fatal("Error: Function doesn't return values at " + fmt.Sprint(node.Token.Line))
	}
	node.Parsed = true
	node.Code = Values.RETURN
	index, t := m.parseExpression(node, 0)
	kind := m.CurrentFunction.Return
	if len(node.Children) <= index || node.Children[index].Type != NEWLINE {
		log.Fatal("Error: Expecting end of line at " + fmt.Sprint(node.Token.Line))
	}
	node.Children[index].Parsed = true
	node.Children[index].Code = ";\n"
	if t != kind {
		log.Fatal("Error: Mismatches types at " + fmt.Sprint(node.Token.Line))
	}
}

func (m *Module) parseFunctionCall(node *Node) {
	f, _ := m.getFunction(node.Token.Value)
	node.Code = node.Token.Value
	node.Parsed = true
	if node.Children[0].Type == LEFT_PAREN {
		if node.Children[1].Type != RIGHT_PAREN {
			log.Fatal("Error: Expecting ')' at " + fmt.Sprint(node.Token.Line))
		}
		node.Children[1].Parsed = true
		node.Children[1].Code = ")"
		node.Children[0].Parsed = true
		node.Children[0].Code = "("
		params := node.Children[0].Children
		childs := len(params)
		l := len(f.Params)
		if (l > 1 && l != (childs+1)/2) || (l == 1 && l != childs) {
			log.Fatal("Error: Wrong size of parameters at " + fmt.Sprint(node.Token.Line))
		}
		for i := 0; i < l; i++ {
			_, t := m.parseExpression(node.Children[0], i*2)
			kind := f.Params[i].Type.Kind
			if t != kind {
				log.Fatal("Error: Mismatches types at " + fmt.Sprint(node.Token.Line))
			}
			if i*2+1 >= childs {
				break
			}
			if params[i*2+1].Type != COMMA {
				log.Fatal("Error: Expecting comma at " + fmt.Sprint(node.Token.Line))
			}
			params[i*2+1].Code = ", "
			params[i*2+1].Parsed = true
		}

		if f.Return == VOID || (len(node.Children) > 2 && node.Children[2].Type == NEWLINE) {
			node.Children[1].Code += ";\n"
		}
	}
}
