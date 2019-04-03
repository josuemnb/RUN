package run

import "os"

var (
	keywords map[string]int
)

func init() {
	keywords = make(map[string]int)
	keywords["and"] = AND
	keywords["else"] = ELSE
	keywords["false"] = FALSE
	keywords["println"] = PRINTLN
	keywords["print"] = PRINT
	keywords["break"] = BREAK
	keywords["if"] = IF
	keywords["nil"] = NIL
	keywords["or"] = OR
	keywords["return"] = RETURN
	keywords["super"] = SUPER
	keywords["this"] = THIS
	keywords["true"] = TRUE
	keywords["loop"] = LOOP
	keywords["main"] = MAIN
	keywords["loop"] = LOOP
	keywords["break"] = BREAK
	keywords["asm"] = ASM
	keywords["cpp"] = CPP
	keywords["import"] = IMPORT
	keywords["module"] = MODULE
}

type Scanner struct {
	start   int
	current int
	line    int
	source  []byte
	length  int
	tokens  []Token
	col     int
	hasMain bool
}

func NewScanner(source []byte, isModule bool) *Scanner {
	s := new(Scanner)
	s.source = source
	s.length = len(source)
	s.tokens = make([]Token, 0)
	s.line = 1
	s.col = 1
	s.hasMain = isModule
	return s
}

func (s *Scanner) HasMain() bool {
	return s.hasMain
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.start = s.current
	s.addTokenVal(EOL, []byte(""))
	s.addTokenVal(EOF, []byte(""))
	return s.tokens
}

func (s *Scanner) addToken(t int) {
	s.addTokenVal(t, nil)
}

func (s *Scanner) addTokenVal(t int, bytes []byte) {
	if bytes == nil {
		s.tokens = append(s.tokens, Token{Type: t, Lexeme: string(s.source[s.start:s.current]), Line: s.line, Col: s.col})
	} else {
		s.tokens = append(s.tokens, Token{Type: t, Lexeme: string(bytes), Line: s.line, Col: s.col})
	}
}

func (s *Scanner) scanToken() {
	b := s.advance()
	switch b {
	case '(':
		s.addToken(LEFT_PAREN)
		s.col++
	case ')':
		s.addToken(RIGHT_PAREN)
		s.col++
	case '{':
		s.addToken(LEFT_BRACE)
		s.col++
	case '}':
		s.addToken(RIGHT_BRACE)
		s.col++
	case '[':
		s.addToken(LEFT_BRACKET)
		s.col++
	case ']':
		s.addToken(RIGHT_BRACKET)
		s.col++
	case ',':
		s.addToken(COMMA)
		s.col++
	case '.':
		if s.match('.') {
			s.addToken(RANGE)
			s.col += 2
		} else {
			s.addToken(DOT)
			s.col++
		}
	case '-':
		if s.match('-') {
			s.addToken(DECREMENT)
			s.col += 2
		} else if s.match('=') {
			s.addToken(MINUS_EQUAL)
			s.col += 2
		} else if s.match('>') {
			s.addToken(INTERFACE)
			s.col += 2
		} else {
			s.addToken(MINUS)
			s.col++
		}
	case '+':
		if s.match('+') {
			s.addToken(INCREMENT)
			s.col += 2
		} else if s.match('=') {
			s.addToken(PLUS_EQUAL)
			s.col += 2
		} else {
			s.addToken(PLUS)
		}
		s.col++
	case ';':
		s.addToken(SEMICOLON)
		s.col++
	case '*':
		s.addToken(STAR)
		s.col++
	case '$':
		s.addMessage()
	case ':':
		s.addToken(DECLARE)
		s.col++
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
			s.col += 2
		} else {
			s.addToken(BANG)
			s.col++
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
			s.col += 2
		} else {
			s.addToken(EQUAL)
			s.col++
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
			s.col += 2
		} else {
			s.addToken(LESS)
			s.col++
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
			s.col += 2
		} else {
			s.addToken(GREATER)
			s.col++
		}
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.peek() != '\n' {
				s.current++
			}
			s.line++
			s.col = 1
		} else {
			s.addToken(SLASH)
			s.col++
		}
	// case ' ', '\t':
	// 	s.col++
	// 	s.addToken(OTHER)
	case '\r':
		return
	case '\n':
		s.line++
		s.addToken(EOL)
		s.col = 1
	case '"', '\'':
		s.addString(b)
	default:
		if s.isDigit(b) {
			for s.isDigit(s.peek()) {
				s.advance()
			}
			p := s.current
			if s.peek() == '.' {
				s.advance()
				if s.peek() == '.' {
					s.addTokenVal(NUMBER, s.source[s.start:p])
					s.addToken(RANGE)
					s.advance()
					return
				}
				for s.isDigit(s.peek()) {
					s.advance()
				}
				s.addTokenVal(REAL, s.source[s.start:s.current])
			} else {
				s.addTokenVal(NUMBER, s.source[s.start:s.current])
			}
			s.col += s.current - s.start
		} else if s.isAlpha(b) {
			for s.isAlphaNum(s.peek()) {
				s.advance()
			}
			id := string(s.source[s.start:s.current])
			if k, ok := keywords[id]; ok {
				s.addToken(k)
				if k == MAIN {
					if s.hasMain {
						println("Main alreay declared")
						os.Exit(-2)
					}
					s.hasMain = true
				}
			} else {
				s.addToken(IDENTIFIER)
			}
			s.col += len(id)
			// } else {
			// 	s.col++
			// 	s.addToken(OTHER)
		}
	}
}

func (s *Scanner) isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (s *Scanner) isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func (s *Scanner) isAlphaNum(b byte) bool {
	return s.isAlpha(b) || s.isDigit(b)
}

func (s *Scanner) addMessage() {
	s.start = s.current
	s.advance()
	for !s.isAtEnd() && s.peek() != '$' {
		if s.peek() == '\n' {
			s.line++
			s.col = 1
		}
		s.advance()
	}
	if s.isAtEnd() {
		println("Error: Missing end of string at", s.line)
		os.Exit(-1)
	}
	s.addTokenVal(MESSAGE, s.source[s.start:s.current])
	s.advance()
	s.col = s.start
}

func (s *Scanner) addString(b byte) {
	s.start = s.current
	// s.advance()
	for !s.isAtEnd() && s.peek() != b {
		if s.peek() == '\n' {
			s.line++
			s.col = 1
		}
		s.advance()
	}
	if s.isAtEnd() {
		println("Error: Missing end of string at", s.line)
		os.Exit(-1)
	}
	s.addTokenVal(QUOTE, s.source[s.start:s.current])
	s.advance()
	s.col = s.start
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= s.length
}

func (s *Scanner) match(b byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != b {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}
