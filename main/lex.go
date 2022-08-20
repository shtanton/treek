package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type stateFunc func(*lexer) stateFunc

type lexer struct {
	input string
	start int
	pos int
	width int
	nestingLevel int
	tokenStream chan Token
}

func (l *lexer) run() {
	for state := lexBlockStart; state != nil; {
		state = state(l)
	}
	close(l.tokenStream)
}

func (l *lexer) emit(t TokenType) {
	l.tokenStream <- Token{
		typ: t,
		val: l.input[l.start:l.pos],
	}
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) stateFunc {
	l.tokenStream <- Token{
		typ: TokenErr,
		val: fmt.Sprintf(format, args...),
	}
	return nil
}

const eof rune = -1

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) reset() {
	l.pos = l.start
}

func (l *lexer) peek() rune {
	w := l.width
	r := l.next()
	l.backup()
	l.width = w
	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptAll(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {}
	l.backup()
}

func (l *lexer) acceptPassing(valid func(rune) bool) bool {
	if valid(l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptAllPassing(valid func(rune) bool) {
	for valid(l.next()) {}
	l.backup()
}

type TokenType int

const (
	TokenErr TokenType = iota // Lexing error
	TokenEOF // end of file
	TokenNumber // number literal
	TokenIdentifier // identifier
	TokenAdd // +
	TokenSub // -
	TokenAst // *
	TokenDiv // /
	TokenDot // .
	TokenComma // ,
	TokenSemicolon // ;
	TokenLParen // (
	TokenRParen // )
	TokenLBrace // {
	TokenRBrace // }
	TokenLBrack // [
	TokenRBrack // ]
	TokenIndexPattern // An index pattern segment
	TokenAssign // =
	TokenCircum // ^
	TokenDoubleQuote // "
	TokenStringLiteral // String literal
	TokenAddAssign // +=
	TokenSubAssign // -=
	TokenAstAssign // *=
	TokenDivAssign // /=
	TokenEqual // ==
	TokenNotEqual // !=
	TokenNot // !
)

type Token struct {
	typ TokenType
	val string
}

func (t Token) String() string {
	switch t.typ {
	case TokenEOF:
		return "EOF"
	case TokenErr:
		return t.val
	}
	if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

func Lex(input string) chan Token {
	l := &lexer{
		input: input,
		tokenStream: make(chan Token),
	}
	go l.run()
	return l.tokenStream
}

const (
	whitespace string = " \t"
	whitespaceNewlines string = " \t\r\n"
)

func isAlpha(r rune) bool {
	return ('a' <= r && r < 'z') || ('A' <= r && r <= 'Z')
}
func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}
func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}
func isIdentifierStartRune(r rune) bool {
	return isAlpha(r) || r == '_' || r == '$'
}
func isIdentifierRune(r rune) bool {
	return isIdentifierStartRune(r) || isDigit(r)
}

func lexBlockStart(l *lexer) stateFunc {
	l.acceptAll(whitespace)
	l.ignore()
	if l.peek() == eof {
		l.emit(TokenEOF)
		return nil
	}
	if l.accept("^") {
		l.emit(TokenCircum)
	}
	if l.accept("{") {
		l.emit(TokenLBrace)
		l.nestingLevel += 1
		return lexStartAction
	}
	return lexPattern
}

func lexPattern(l *lexer) stateFunc {
	r := l.next()
	switch {
		case isIdentifierRune(r):
			l.backup()
			return lexIndexPattern
		case r == '(':
			l.nestingLevel += 1
			l.emit(TokenLParen)
			return lexFilterPattern
		case r == '*':
			l.emit(TokenAst)
			return lexPatternEnd
	}
	return l.errorf("Invalid Pattern")
}

// Next rune is identifier rune
func lexIndexPattern(l *lexer) stateFunc {
	l.acceptAllPassing(isIdentifierRune)
	l.emit(TokenIndexPattern)
	return lexPatternEnd
}

// Previous rune is (
func lexFilterPattern(l *lexer) stateFunc {
	for state := lexAction; state != nil; {
		state = state(l)
	}
	return lexPatternEnd
}

func lexPatternEnd(l *lexer) stateFunc {
	if l.accept(".") {
		l.emit(TokenDot)
		return lexPattern
	}
	l.acceptAll(whitespaceNewlines)
	l.ignore()
	if !l.accept("{") {
		if l.peek() == eof {
			l.emit(TokenEOF)
			return nil
		}
		return l.errorf("Missing Action")
	}
	l.emit(TokenLBrace)
	l.nestingLevel += 1
	return lexStartAction
}

// Previous rune is {
func lexStartAction(l *lexer) stateFunc {
	for state := lexAction; state != nil; {
		state = state(l)
	}
	return lexBlockStart
}

func lexAction(l *lexer) stateFunc {
	l.acceptAll(whitespaceNewlines)
	l.ignore()
	doubleCharTokens := map[rune]map[rune]TokenType{
		'+': {
			'=': TokenAddAssign,
		},
		'-': {
			'=': TokenSubAssign,
		},
		'*': {
			'=': TokenAstAssign,
		},
		'/': {
			'=': TokenDivAssign,
		},
		'=': {
			'=': TokenEqual,
		},
		'!': {
			'=': TokenNotEqual,
		},
	}
	charTokens := map[rune]TokenType{
		'+': TokenAdd,
		'-': TokenSub,
		'/': TokenDiv,
		'*': TokenAst,
		'.': TokenDot,
		',': TokenComma,
		';': TokenSemicolon,
		'=': TokenAssign,
		'!': TokenNot,
	}
	r := l.next()
	charToken, isCharToken := charTokens[r]
	doubleCharMap, hasDoubleCharMap := doubleCharTokens[r]
	if hasDoubleCharMap {
		doubleCharToken, hasDoubleCharToken := doubleCharMap[l.next()]
		if hasDoubleCharToken {
			l.emit(doubleCharToken)
			return lexAction
		} else {
			l.backup()
		}
	}
	switch {
		case r == eof:
			return l.errorf("Unclosed Action")
		case r == '(':
			l.nestingLevel += 1
			l.emit(TokenLParen)
			return lexAction
		case r == ')':
			l.nestingLevel -= 1
			l.emit(TokenRParen)
			if l.nestingLevel > 0 {
				return lexAction
			}
			return nil
		case r == '{':
			l.nestingLevel += 1
			l.emit(TokenLBrace)
			return lexAction
		case r == '}':
			l.nestingLevel -= 1
			l.emit(TokenRBrace)
			if l.nestingLevel > 0 {
				return lexAction
			}
			return nil
		case r == '[':
			l.nestingLevel += 1
			l.emit(TokenLBrack)
			return lexAction
		case r == ']':
			l.nestingLevel -= 1
			l.emit(TokenRBrack)
			if l.nestingLevel > 0 {
				return lexAction
			}
			return nil
		case isCharToken:
			l.emit(charToken)
			return lexAction
		case isDigit(r):
			return lexNumberLiteral
		case isIdentifierStartRune(r):
			return lexIdentifier
		case r == '"':
			l.emit(TokenDoubleQuote)
			return lexStringLiteral
	}
	return l.errorf("Invalid Token: " + string(r))
}

// Just accepted the first digit
func lexNumberLiteral(l *lexer) stateFunc {
	l.acceptAllPassing(isDigit)
	if l.accept(".") {
		l.acceptAllPassing(isDigit)
	}
	if isIdentifierRune(l.peek()) {
		l.next()
		return l.errorf("Bad number: %q", l.input[l.start:l.pos])
	}
	l.emit(TokenNumber)
	return lexAction
}

func lexStringLiteral(l *lexer) stateFunc {
	for {
		switch l.next() {
			case eof:
				return l.errorf("Missing closing quote in string literal")
			case '"':
				l.backup()
				l.emit(TokenStringLiteral)
				l.next()
				l.emit(TokenDoubleQuote)
				return lexAction
		}
	}
}

// Just accepted the first character
func lexIdentifier(l *lexer) stateFunc {
	l.acceptAllPassing(isIdentifierRune)
	l.emit(TokenIdentifier)
	return lexAction
}
