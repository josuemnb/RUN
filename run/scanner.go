package run

import (
	"container/list"
	"fmt"
	"log"
)

var (
	keywords map[string]int
	builtin  map[string]int
	queue    *list.List
)

type QueueNode struct {
	Type   int
	Node   *Node
	Return bool
}

func init() {
	keywords = make(map[string]int)
	keywords["and"] = AND
	keywords["else"] = ELSE
	keywords["break"] = BREAK
	keywords["if"] = IF
	keywords["nil"] = NIL
	keywords["or"] = OR
	keywords["return"] = RETURN
	keywords["super"] = SUPER
	keywords["this"] = THIS
	keywords["loop"] = LOOP
	keywords["main"] = MAIN
	keywords["loop"] = LOOP
	keywords["break"] = BREAK
	keywords["asm"] = ASM
	keywords["C"] = C
	keywords["import"] = IMPORT
	keywords["new"] = NEW

	builtin = make(map[string]int)
	builtin["println"] = PRINTLN
	builtin["print"] = PRINT

	queue = list.New()
}

type Scanner struct {
	Root     *Node
	Bytes    []byte
	pointer  int
	current  *Node
	lastNode *Node
	HasMain  bool
	line     int
	col      int
	last     int
	size     int
}

func NewScanner(bytes []byte) *Scanner {
	p := new(Scanner)
	p.Bytes = bytes
	p.size = len(bytes)
	p.Root = NewNode()
	p.current = p.Root
	p.line = 1
	p.last = 0
	p.col = 1
	return p
}

func (p *Scanner) printTree() {
	printTree(p.Root, 0)
}

func printTreeCode(n *Node, tab int) {
	if n.Type == NEWLINE {
		return
	}
	for i := 0; i < tab; i++ {
		print(" ")
	}
	if n.Code != "" && n.Code != "\n" {
		print(tab, "═> ", n.Code, "\n")
	}
	// print("╚═> ", n.Token.Value, "\n")
	for _, v := range n.Children {
		printTreeCode(v, tab+1)
	}
}

func printTree(n *Node, tab int) {
	if n.Type == NEWLINE {
		return
	}
	for i := 0; i < tab; i++ {
		print(" ")
	}
	print(tab, "═> ", n.Token.Value, "\n")
	// print("╚═> ", n.Token.Value, "\n")
	for _, v := range n.Children {
		printTree(v, tab+1)
	}
}

func (p *Scanner) scan() {
	for p.pointer < p.size {
		for p.pointer < p.size && (p.Bytes[p.pointer] == ' ' || p.Bytes[p.pointer] == '\t') {
			p.pointer++
			p.col++
		}
		if p.pointer >= p.size {
			break
		}
		node := NewNode()
		node.Parent = p.current
		node.Index = len(p.current.Children)
		add := false
		ret := false
		back := false
		p.last = p.pointer
		if p.isDigit(p.Bytes[p.pointer]) {
			p.last = p.pointer
			for p.pointer < p.size && p.isDigit(p.Bytes[p.pointer]) {
				p.pointer++
				p.col++
			}
			if p.Bytes[p.pointer] == '.' {
				p.pointer++
				if p.Bytes[p.pointer] == '.' {
					node.Type = NUMBER
					node.Token = Token{Value: string(p.Bytes[p.last : p.pointer-1]), Line: p.line, Col: p.col}
					p.current.Children = append(p.current.Children, node)
					p.pointer--
					p.col--
					continue
				} else {
					p.col++
					for p.isDigit(p.Bytes[p.pointer]) {
						p.pointer++
						p.col++
					}
					node.Type = REAL
				}
			} else {
				node.Type = NUMBER
			}
			node.Token = Token{Value: string(p.Bytes[p.last:p.pointer]), Type: node.Type, Line: p.line, Col: p.col}
			p.current.Children = append(p.current.Children, node)
			p.lastNode = node
			continue
		} else if p.isAlpha(p.Bytes[p.pointer]) {
			p.last = p.pointer
			for p.pointer < p.size && p.isAlphaNum(p.Bytes[p.pointer]) {
				p.pointer++
				p.col++
			}
			id := string(p.Bytes[p.last:p.pointer])
			if key, ok := keywords[id]; ok {
				if key == MAIN {
					p.HasMain = true
				}
				node.Type = KEYWORD
			} else if id == "false" || id == "true" {
				node.Type = BOOL
			} else if _, ok := builtin[id]; ok {
				node.Type = BUILTIN
			} else {
				node.Type = IDENTIFIER
			}
			node.Token = Token{Value: id, Type: node.Type, Line: p.line, Col: p.col}
			p.current.Children = append(p.current.Children, node)
			if node.Type == KEYWORD || node.Type == BUILTIN {
				p.current = node
			}
			p.lastNode = node
			continue
		}
		switch p.Bytes[p.pointer] {
		case '\n':
			p.col = -1
			p.line++
			node.Type = NEWLINE
			add = true
		case ')':
			q := queue.Front()
			n := q.Value.(*QueueNode)
			if n.Type != LEFT_PAREN {
				log.Fatal("Error: Parentesis out of place at " + fmt.Sprint(p.line))
			}
			if n.Return {
				back = true
			}
			queue.Remove(q)
			p.current = n.Node
			ret = true
			add = true
			node.Type = RIGHT_PAREN
		case '(':
			r := false
			l := len(p.current.Children)
			if p.current.Type == IDENTIFIER || p.current.Type == KEYWORD {
				p.current = p.lastNode
				node.Parent = p.current
				r = true
			} else if l > 0 && (p.current.Children[l-1].Type == IDENTIFIER || p.current.Children[l-1].Type == KEYWORD) {
				// r = true
				p.current = p.current.Children[l-1]
				node.Parent = p.current
			}
			queue.PushFront(&QueueNode{Node: node, Type: LEFT_PAREN, Return: r})
			p.addNode(LEFT_PAREN)
		case '}':
			q := queue.Front()
			n := q.Value.(*QueueNode)
			if n.Type != LEFT_BRACE {
				log.Fatal("Error: Parentesis out of place at " + fmt.Sprint(p.line))
			}
			p.current = n.Node
			if n.Return {
				back = true
			}
			queue.Remove(q)
			ret = true
			add = true
			node.Type = RIGHT_BRACE
		case '{':
			r := false
			l := len(p.current.Children)
			if p.current.Type == KEYWORD || p.current.Type == IDENTIFIER {
				r = true
				//TODO Confirmar se o seguinte bloco aplica-se em qualquer teste
			} else if l > 0 && (p.current.Children[l-1].Type == IDENTIFIER || p.current.Children[l-1].Type == KEYWORD) {
				r = true
				p.current = p.current.Children[l-1]
				node.Parent = p.current
			}
			queue.PushFront(&QueueNode{Node: node, Type: LEFT_BRACE, Return: r})
			p.addNode(LEFT_BRACE)
		case ']':
			q := queue.Front()
			n := q.Value.(*QueueNode)
			if n.Type != LEFT_BRACKET {
				log.Fatal("Error: Brackets out of place at " + fmt.Sprint(p.line))
			}
			p.current = n.Node
			queue.Remove(q)
			ret = true
			add = true
			node.Type = RIGHT_BRACKET
		case '[':
			r := false
			l := len(p.current.Children)
			if p.current.Type == KEYWORD || p.current.Type == IDENTIFIER {
				r = true
				//TODO Confirmar se o seguinte bloco aplica-se em qualquer teste
			} else if l > 0 && (p.current.Children[l-1].Type == IDENTIFIER || p.current.Children[l-1].Type == KEYWORD) {
				r = true
				p.current = p.current.Children[l-1]
				node.Parent = p.current
			}
			queue.PushFront(&QueueNode{Node: node, Type: LEFT_BRACKET, Return: r})
			p.addNode(LEFT_BRACKET)
		case '+':
			if p.Bytes[p.pointer+1] == '+' {
				node.Type = INCREMENT
				p.pointer++
				p.col++
			} else if p.Bytes[p.pointer+1] == '=' {
				node.Type = PLUS_EQUAL
				p.pointer++
				p.col++
			} else {
				node.Type = PLUS
			}
			add = true
		case '-':
			if p.Bytes[p.pointer+1] == '-' {
				node.Type = DECREMENT
				p.pointer++
				p.col++
			} else if p.Bytes[p.pointer+1] == '=' {
				node.Type = MINUS_EQUAL
				p.pointer++
				p.col++
			} else {
				node.Type = MINUS
			}
			add = true
		case '*':
			if p.Bytes[p.pointer+1] == '=' {
				node.Type = MUL_EQUAL
				p.pointer++
				p.col++
			} else {
				node.Type = MUL
			}
			add = true
		case '/':
			if p.Bytes[p.pointer+1] == '=' {
				node.Type = DIV_EQUAL
				p.pointer++
				p.col++
			} else {
				node.Type = DIV
			}
			add = true
		case '%':
			node.Type = MOD
			add = true
		case '!':
			node.Type = NOT
			add = true
		case '&':
			if p.Bytes[p.pointer+1] == '&' {
				node.Type = AND
				p.pointer++
				p.col++
			} else {
				node.Type = BIT_AND
			}
			add = true
		case '|':
			if p.Bytes[p.pointer+1] == '|' {
				node.Type = OR
				p.pointer++
				p.col++
			} else {
				node.Type = BIT_OR
			}
			add = true
		case '^':
		case '~':
		case '\\':
			p.pointer++
			p.last = p.pointer
			for p.pointer < p.size && p.Bytes[p.pointer] != '\n' {
				p.pointer++
			}
		case '\'':
			p.pointer++
			for p.pointer < p.size && p.Bytes[p.pointer] != '\'' {
				p.pointer++
				p.col++
			}
			add = true
			node.Type = QUOTE
		case '"':
			p.pointer++
			for p.pointer < p.size && p.Bytes[p.pointer] != '"' {
				p.pointer++
				p.col++
			}
			add = true
			node.Type = QUOTE
		case ':':
			node.Type = DECLARE
			add = true
		case '.':
			if p.Bytes[p.pointer+1] == '.' {
				// p.current = p.current.Parent
				node.Type = RANGE
				p.pointer++
				p.col++
			} else {
				node.Type = DOT
			}
			add = true
		case ',':
			add = true
			node.Type = COMMA
		case '>':
			if p.Bytes[p.pointer+1] == '=' {
				node.Type = GREATER_EQUAL
				p.pointer++
				p.col++
			} else {
				node.Type = GREATER
			}
			add = true
		case '<':
			if p.Bytes[p.pointer+1] == '=' {
				node.Type = LESS_EQUAL
				p.pointer++
				p.col++
			} else {
				node.Type = LESS
			}
			add = true
		case '=':
			if p.Bytes[p.pointer+1] == '=' {
				node.Type = EQUAL_EQUAL
				p.pointer++
				p.col++
			} else {
				node.Type = EQUAL
				// change = true
			}
			add = true
		}
		if ret {
			p.current = p.current.Parent
			p.lastNode = p.current
		}
		if add {
			node.Token = Token{Value: string(p.Bytes[p.last : p.pointer+1]), Type: node.Type, Line: p.line, Col: p.col}
			p.current.Children = append(p.current.Children, node)
			p.lastNode = node
		}
		if back {
			p.current = p.current.Parent
			p.lastNode = p.current
		}
		p.pointer++
		p.col++
	}
	if queue.Len() > 0 {
		log.Fatal("Error: Missing parentesis or block or brackets closing")
		// os.Exit(-3)
	}
	// if p.current != p.Root {
	// 	println("Error: Unexpected end of file")
	// 	os.Exit(-1)
	// }
}

func (p *Scanner) addNode(kind int) {
	temp := NewNode()
	temp.Index = len(p.current.Children)
	temp.Parent = p.current
	temp.Type = kind
	temp.Token = Token{Value: string(p.Bytes[p.pointer]), Type: kind, Line: p.line, Col: p.col}
	p.current.Children = append(p.current.Children, temp)
	p.current = temp
	p.lastNode = temp
}

func (p *Scanner) addNodeValue(kind int, s string) {
	temp := NewNode()
	temp.Index = len(p.current.Children)
	temp.Parent = p.current
	temp.Type = kind
	temp.Token = Token{Value: s, Type: kind, Line: p.line, Col: p.col}
	p.current.Children = append(p.current.Children, temp)
	p.current = temp
	p.lastNode = temp
}

func (p *Scanner) isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (p *Scanner) isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func (p *Scanner) isAlphaNum(b byte) bool {
	return p.isAlpha(b) || p.isDigit(b)
}
