package run

type Type struct {
	Name        string
	Kind        int
	Class       Class
	Collection  int
	Used        int
	Real        string
	Parent      *Module
	IsInterface bool
	Interface   Interface
}

type Protection int

const (
	PUBLIC Protection = iota
	PROTECTED
	PRIVATE
)

var (
	// kinds map[string]Kind
	// types   []Type
	typeIdx int
	// typesUsed   map[string]bool
	collections map[string]*Type
)

func init() {
	// typesUsed = make(map[string]bool)
	collections = make(map[string]*Type)
	typeIdx = STRING + 1
	// types = make([]Type, 0)
}

func (m *Module) updateType(t *Type) {
	m.Types[t.Name] = t
}

func (m *Module) addType(t *Type) {
	if _, ok := m.Types[t.Name]; ok {
		m.error("Type already defined", 1)
	}
	if t.Kind == 0 {
		t.Kind = typeIdx
		typeIdx++
	}
	if m.Type == MODULE {
		t.Real = m.Name + "_"
	}
	if t.Kind >= STRING && t.IsInterface == false && t.Collection == 0 {
		t.Real += "class_"
	}
	t.Parent = m
	t.Real += t.Name
	m.Types[t.Name] = t
}

func (m *Module) getTypeByName(n string) *Type {
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

func (m *Module) getTypeByKind(i int) *Type {
	for _, t := range m.Types {
		if t.Kind == i {
			return t
		}
	}
	return nil
}

func (m *Module) typeOf(node Node) Type {
	switch node.Type {
	case LITERAL:
		l := node.Value.(Literal)
		return l.Type
	case NEWLINE:
		return Type{Kind: NEWLINE}
	case IDENTIFIER:
		i := node.Value.(Identifier)
		return m.getType(i.Name)
	case UNARY:
		u := node.Value.(Unary)
		return m.typeOf(u.Right)
	case BINARY:
		b := node.Value.(Binary)
		t0 := m.typeOf(b.Left)
		t1 := m.typeOf(b.Right)
		if t1.Kind == 0 {
			return t0
		}
		return t1
	case CALL:
		c := node.Value.(Call)
		if c.Kind == CLASS || c.Kind == THIS {
			return *m.getTypeByName(c.Name)
		} else if c.Kind == FUNCTION {
			return c.Return
		} else if c.Kind == METHOD {
			return c.Return
		}
	case BRACKETS:
		b := node.Value.(Bracket)
		t := m.getType(b.Left.Name)
		if c, ok := collections[t.Name]; ok {
			return c.Class.Methods["get"].Return
		}
		return t
	case GROUPING:
		return m.typeOf(node.Value.(Grouping).Group)
	case DOT:
		d := node.Value.(Dot)
		return Type{Kind: DOT, Class: m.getClass(m.typeOf(d.Left).Class.Name)}
	case NULL:
		return Type{Kind: NULL}
	}
	return Type{}
}

func typeToRepr(t int) string {
	switch t {
	case NUMBER:
		return "%lld"
	case REAL:
		return "%f"
	case STRING, NEWLINE, QUOTE, BOOL:
		return "%s"
	}
	return "%x"
}

type Param struct {
	Name       string
	Type       Type
	Protection int
}
