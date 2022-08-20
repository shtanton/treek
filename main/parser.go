package main

import (
	"strconv"
	"fmt"
)

type InstructionBasic int
const (
	InstructionAdd InstructionBasic = iota
	InstructionSub
	InstructionDiv
	InstructionMul
	InstructionIgnore
	InstructionPushNull
	InstructionAssign
	InstructionIndex
	InstructionDup
	InstructionEqual
	InstructionNot
)

type InstructionPushNumber float64
type InstructionPushVariable string
type InstructionPushString string

type Subroutine int
const (
	SubroutinePrintln Subroutine = iota
)
type InstructionCall struct {
	subroutine Subroutine
	nargs int
}

type Instruction interface {
	debug()
	eval(*EvalState)
}

func (i InstructionBasic) debug() {
	switch i {
		case InstructionAdd:
			fmt.Println("Add")
		case InstructionSub:
			fmt.Println("Sub")
		case InstructionDiv:
			fmt.Println("Div")
		case InstructionMul:
			fmt.Println("Mul")
		case InstructionIgnore:
			fmt.Println("Ignore")
		case InstructionPushNull:
			fmt.Println("Push Null")
		case InstructionAssign:
			fmt.Println("Assign")
		case InstructionIndex:
			fmt.Println("Index")
		case InstructionDup:
			fmt.Println("Dup")
		case InstructionEqual:
			fmt.Println("Equal")
		case InstructionNot:
			fmt.Println("Not")
		default:
			fmt.Println("Unknown Basic Instruction")
	}
}

func (n InstructionPushNumber) debug() {
	fmt.Printf("Push Number: %v\n", n)
}

func (i InstructionCall) debug() {
	subroutines := map[Subroutine]string {
		SubroutinePrintln: "println",
	}
	fmt.Printf("Calling %v with %v arguments\n", subroutines[i.subroutine], i.nargs)
}

func (s InstructionPushVariable) debug() {
	fmt.Printf("Push variable: %v\n", s)
}

func (s InstructionPushString) debug() {
	fmt.Printf("Push string: %q\n", s)
}

type Expression []Instruction

type PatternSegmentIndex string
type PatternSegmentFilter Expression
type PatternSegmentBasic int
const (
	PatternSegmentAll PatternSegmentBasic = iota
)

type PatternSegment interface {
	debug()
	matches(*EvalState, []TreePathSegment, TreePathSegment) bool
}

type Pattern struct {
	segments []PatternSegment
	isFirst bool
}

func (s PatternSegmentIndex) debug() {
	fmt.Printf("Index: %q\n", s)
}

func (s PatternSegmentFilter) debug() {
	fmt.Println("Filter: (")
	for _, instruction := range s {
		instruction.debug()
	}
	fmt.Println(")")
}

func (s PatternSegmentBasic) debug() {
	switch s {
		case PatternSegmentAll:
			fmt.Println("All")
		default:
			panic("Invalid basic pattern segment")
	}
}

type Block struct {
	pattern Pattern
	action Expression
}

type Program struct {
	blocks []Block
}

func (p Program) debug() {
	for _, block := range p.blocks {
		fmt.Println("\nPattern:")
		if block.pattern.isFirst {
			fmt.Println("First")
		} else {
			fmt.Println("Last")
		}
		for _, s := range block.pattern.segments {
			s.debug()
		}
		fmt.Println("\nAction:")
		for _, a := range block.action {
			a.debug()
		}
	}
}

type parser struct {
	tokenStream chan Token
	prevToken Token
	wasRewound bool
}

func (p *parser) next() Token {
	if p.wasRewound {
		p.wasRewound = false
		return p.prevToken
	}
	p.prevToken = <-p.tokenStream
	if p.prevToken.typ == TokenErr {
		fmt.Printf("Error: %q\n", p.prevToken.val)
		panic("Lexing error")
	}
	return p.prevToken
}

func (p *parser) rewind() {
	p.wasRewound = true
}

func (p *parser) accept(typ TokenType) (string, bool) {
	token := p.next()
	if typ == token.typ {
		return token.val, true
	}
	p.rewind()
	return "", false
}

func (p *parser) peek() Token {
	token := p.next()
	p.rewind()
	return token
}

func (p *parser) parsePatternSegment() (segment PatternSegment, action bool, eof bool) {
	token := p.next()
	switch token.typ {
		case TokenEOF:
			p.rewind()
			return nil, false, true
		case TokenLBrace:
			p.rewind()
			return nil, true, false
		case TokenIndexPattern:
			return PatternSegmentIndex(token.val), false, false
		case TokenLParen:
			filter, noExpression := p.parseExpression(0)
			if noExpression {
				panic("Missing expression in filter")
			}
			_, hasCloseParen := p.accept(TokenRParen)
			if !hasCloseParen {
				panic("Missing close paren")
			}
			return PatternSegmentFilter(filter), false, false
		case TokenAst:
			return PatternSegmentAll, false, false
		default:
			panic("Expected pattern segment")
	}
}

func (p *parser) parsePattern() (pattern Pattern, eof bool) {
	_, pattern.isFirst = p.accept(TokenCircum)
	segment, action, eof := p.parsePatternSegment()
	if eof {
		return pattern, true
	} else if action {
		return pattern, false
	}
	pattern.segments = append(pattern.segments, segment)
	for {
		_, hasAnotherSegment := p.accept(TokenDot)
		if !hasAnotherSegment {
			break
		}
		segment, action, eof := p.parsePatternSegment()
		if eof || action {
			panic("Expected pattern segment")
		}
		pattern.segments = append(pattern.segments, segment)
	}
	return pattern, false
}

func (p *parser) parseExpression(minPower int) (expr Expression, noExpression bool) {
	token := p.next()
	switch token.typ {
		case TokenEOF:
			p.rewind()
			return nil, true
		case TokenNot:
			e, noExpression := p.parseExpression(14)
			if noExpression {
				panic("Missing expression after !")
			}
			expr = append(expr, e...)
			expr = append(expr, InstructionNot)
		case TokenNumber:
			num, err := strconv.ParseFloat(token.val, 64)
			if err != nil {
				panic("Invalid number")
			}
			expr = append(expr, InstructionPushNumber(num))
		case TokenDoubleQuote:
			s, isStringLiteral := p.accept(TokenStringLiteral)
			if !isStringLiteral {
				panic("Missing string literal")
			}
			_, stringLiteralClosed := p.accept(TokenDoubleQuote)
			if !stringLiteralClosed {
				panic("Missing closing quote for string literal")
			}
			expr = append(expr, InstructionPushString(s))
		case TokenIdentifier:
			_, hasLParen := p.accept(TokenLParen)
			if hasLParen {
				subroutines := map[string]Subroutine {
					"println": SubroutinePrintln,
				}
				subroutine, isSubroutine := subroutines[token.val]
				if !isSubroutine {
					panic("Invalid subroutine")
				}
				nargs := 0
				for {
					e, noExpression := p.parseExpression(0)
					if noExpression {
						break
					}
					expr = append(expr, e...)
					nargs += 1
					_, hasComma := p.accept(TokenComma)
					if !hasComma {
						break
					}
				}
				_, hasRParen := p.accept(TokenRParen)
				if !hasRParen {
					panic("Missing ) for subroutine call")
				}
				expr = append(expr, InstructionCall {subroutine, nargs})
			} else {
				expr = append(expr, InstructionPushVariable(token.val))
			}
		case TokenLParen:
			e, noExpression := p.parseExpression(0)
			if noExpression {
				panic("Missing expression in ()")
			}
			_, hasCloseParen := p.accept(TokenRParen)
			if !hasCloseParen {
				panic("Missing ) in expression")
			}
			expr = append(expr, e...)
		default:
			p.rewind()
			return nil, true
	}
	
	oploop: for {
		token := p.next()
		binops := map[TokenType] struct{
			op InstructionBasic
			left, right int
		} {
			TokenAdd: {InstructionAdd, 10, 11},
			TokenSub: {InstructionSub, 10, 11},
			TokenAst: {InstructionMul, 12, 13},
			TokenDiv: {InstructionDiv, 12, 13},
			TokenAssign: {InstructionAssign, 3, 2},
			TokenEqual: {InstructionEqual, 8, 9},
		}
		binop, isBinop := binops[token.typ]
		assigns := map[TokenType]InstructionBasic {
			TokenAddAssign: InstructionAdd,
			TokenSubAssign: InstructionSub,
			TokenAstAssign: InstructionMul,
			TokenDivAssign: InstructionDiv,
		}
		assignInstruction, isAssign := assigns[token.typ]
		switch {
			case isBinop && binop.left >= minPower:
				e, noExpression := p.parseExpression(binop.right)
				if noExpression {
					panic("Missing expression after operator")
				}
				expr = append(expr, e...)
				expr = append(expr, binop.op)
			case isAssign && 3 >= minPower:
				expr = append(expr, InstructionDup)
				e, noExpression := p.parseExpression(2)
				if noExpression {
					panic("Missing expression after operator")
				}
				expr = append(expr, e...)
				expr = append(expr, assignInstruction, InstructionAssign)
			case token.typ == TokenSemicolon && 0 >= minPower:
				e, noExpression := p.parseExpression(1)
				expr = append(expr, InstructionIgnore)
				if noExpression {
					expr = append(expr, InstructionPushNull)
				} else {
					expr = append(expr, e...)
				}
			case token.typ == TokenDot && 20 >= minPower:
				index, hasIndex := p.accept(TokenIdentifier)
				if !hasIndex {
					panic("Expected identifier after .")
				}
				expr = append(expr, InstructionPushString(index), InstructionIndex)
			case token.typ == TokenNotEqual && 8 >= minPower:
				e, noExpression := p.parseExpression(9)
				if noExpression {
					panic("Missing expression after operator")
				}
				expr = append(expr, e...)
				expr = append(expr, InstructionEqual, InstructionNot)
			default:
				p.rewind()
				break oploop
		}
	}
	
	return expr, false
}

func Parse(tokenStream chan Token) Program {
	p := parser {
		tokenStream: tokenStream,
		wasRewound: false,
	}
	var blocks []Block
	for {
		pattern, eof := p.parsePattern()
		if eof {
			break
		}
		_, hasAction := p.accept(TokenLBrace)
		var action Expression
		if hasAction {
			var noAction bool
			action, noAction = p.parseExpression(0)
			if noAction {
				action = nil
			}
			_, hasActionClose := p.accept(TokenRBrace)
			if !hasActionClose {
				panic("Error: Missing } at end of action")
			}
		}
		blocks = append(blocks, Block {
			pattern: pattern,
			action: action,
		})
	}
	return Program {
		blocks: blocks,
	}
}
