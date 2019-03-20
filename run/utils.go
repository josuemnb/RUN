package run

// func typeOf(i interface{}) reflect.Type {
// 	return reflect.TypeOf(i)
// }

// func valueOf(e Expr) string {
// 	switch e.Type() {
// 	case VARIABLE:
// 		val := e.(Variable)
// 		return "var_" + val.Name
// 	case LITERAL:
// 		val := e.(Literal)
// 		return val.Value
// 	case UNARY:
// 		val := e.(Unary)
// 		return valueOf(val.Right) + val.Operator.Lexeme
// 	case BINARY, LOGICAL:
// 		val := e.(Binary)
// 		l := valueOf(val.Left)
// 		r := valueOf(val.Right)
// 		return l + val.Operator.Lexeme + r
// 		// case NAME:
// 		// 	val := e.(Name)
// 		// 	return
// 	}
// 	return ""
// }

// func valueOfs(e Expr) (v string) {
// 	switch e.Type() {
// 	case VARIABLE:
// 		val := e.(Variable)
// 		v = val.Name
// 	case LITERAL, BOOLEAN:
// 		val := e.(Literal)
// 		v = val.Value
// 	case BINARY:
// 		val := e.(Binary)
// 		lv := valueOf(val.Left)
// 		rv := valueOf(val.Right)
// 		if lv != "" {
// 			v += lv
// 		}
// 		if rv != "" {
// 			v += rv
// 		}
// 		return
// 	case LOGICAL:
// 		val := e.(Logical)
// 		lv := valueOf(val.Left)
// 		rv := valueOf(val.Right)
// 		if lv != "" {
// 			v += lv
// 		}
// 		if rv != "" {
// 			v += rv
// 		}
// 	}
// 	return
// }

// func exprToString(e Expr) (t, s string) {
// 	switch e.Type() {
// 	case VARIABLE:
// 		val := e.(Variable)
// 		s = val.Name
// 		t = val.Kind.Name
// 	case LITERAL:
// 		val := e.(Literal)
// 		s = val.Value
// 		t = tokenTypeString(TokenType(val.kind))
// 	case BOOLEAN:
// 		val := e.(Literal)
// 		s = val.Value
// 		t = "bool"
// 	case BINARY:
// 		val := e.(Binary)
// 		s = valueOf(val.Left) + val.Operator.Lexeme + valueOf(val.Right)
// 	}
// 	return
// }
// func tokenTypeRepr(t TokenType) string {
// 	switch t {
// 	case NUMBER:
// 		return "%d"
// 	case STRING:
// 		return "%s"
// 	}
// 	return "%x"
// }

// func tokenTypeString(t TokenType) string {
// 	switch t {
// 	case NUMBER:
// 		return "int"
// 	case STRING:
// 		return "string"
// 	}
// 	return "string"
// }

// func typeToRepr(t int) string {
// 	switch t {
// 	case 1:
// 		return "%d"
// 	case 2:
// 		return "%s"
// 	}
// 	return "%x"
// }

// func isInt(i interface{}) bool {
// 	return typeOf(i).Kind() == reflect.Int
// }

// func isString(i interface{}) bool {
// 	return typeOf(i).Kind() == reflect.String
// }
