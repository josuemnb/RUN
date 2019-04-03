package run

type Scope struct {
	Variables map[string]Variable
	// Parent    *Scope
}

type Variable struct {
	Type       Type
	Name       string
	Protection Protection
	Array      int
}

func (m *Module) BeginScope() {
	// var scope Scope
	scope := make(map[string]Variable)
	// scope.Parent = &m.Scope
	m.CurScope++
	if len(m.Scopes) <= m.CurScope {
		m.Scopes = append(m.Scopes, scope)
	} else {
		m.Scopes[m.CurScope] = scope
	}
}

func (m *Module) EndScope() {
	m.CurScope--
	// m.Scope = *m.Scope.Parent
}

func (m *Module) exists(n string) bool {
	_, ok := m.Scopes[m.CurScope][n]
	return ok
}

func (m *Module) getType(n string) Type {
	current := m.CurScope
	for {
		if cls, ok := m.Scopes[current][n]; ok {
			return cls.Type
		}
		if current == 0 {
			break
		}
		current--
	}
	for _, m := range m.Modules {
		if t := m.getType(n); t.Kind > 0 {
			t.Class.Name = m.Name + "_" + t.Class.Name
			return t
		}
	}
	return Type{}
}

func (m *Module) getFuncType(v string) Type {
	if t, ok := m.Functions[v]; ok {
		return t.Return
	}
	return Type{}
}
