package run

const (
	VOID int = 1 + iota
	NULL
	VARIABLE
	GROUPING
	CALL
	GET
	SET
	LITERAL
	UNARY
	BINARY
	ASSIGN
	LOGICAL
	BOOLEAN
	NAME
	FUNCTION
	INTERFACE
	EMPTY
	OPERATOR
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_BRACKET
	RIGHT_BRACKET
	BRACKETS
	BETWEEN
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	DECLARE
	EOL
	RANGE
	BANG
	BANG_EQUAL
	PLUS_EQUAL
	MINUS_EQUAL
	MUL
	MUL_MUL
	MUL_EQUAL
	BIT_AND
	BIT_OR
	DIV
	DIV_EQUAL
	EQUAL
	MOD
	NOT
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	KEYWORD
	IDENTIFIER
	PRIVATE
	PUBLIC
	PROTECTED
	AND
	ELSE
	FALSE
	IF
	ELSEIF
	NIL
	OR
	BREAK
	PRINT
	PRINTLN
	RETURN
	SUPER
	EXTENDS
	THIS
	TRUE
	LOOP
	EOF
	CARDINAL
	INCREMENT
	DECREMENT
	EXPRESSION
	PARENTESIS
	MAIN
	BUILTIN
	CLASS
	METHOD
	FIELD
	BLOCK
	PARAM
	SEQUENCE
	NEWLINE
	ASM
	C
	NEW
	IMPORT
	MODULE
	MESSAGE
	QUOTE
	MAP
	LIST
	STACK
	ARRAY
	NUMBER
	BYTE
	REAL
	BOOL
	ROOT
	STRING int = 200
)

type Node struct {
	Type     int
	Token    Token
	Parent   *Node
	Children []*Node
	Code     string
	Parsed   bool
	Index    int
	NotPrint bool
}

func NewNode() *Node {
	n := new(Node)
	n.Children = make([]*Node, 0)
	return n
}

type Type struct {
	Name       string
	Kind       int
	Class      *Class
	Collection *Collection
	// Used        int
	Real   string
	Module *Module
	// IsInterface bool
	// Interface   Interface
}

type Param struct {
	Name       string
	Type       Type
	Protection int
}

type Token struct {
	Type  int
	Value string
	Line  int
	Col   int
}
