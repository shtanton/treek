package main

import (
	"fmt"
	"os"
	"bufio"
)

type TreePathSegment interface{}
type TreeData struct {
	path []TreePathSegment
	value Value
}
type TreeStream chan TreeData

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing program arg")
		return
	}
	input := os.Args[1]
	tokens := Lex(input)
	program := Parse(tokens)
	
	stdin := bufio.NewReader(os.Stdin)
	data := Json(stdin)
	
	Eval(program, data)
}
