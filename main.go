package main

import (
	"fmt"
)

func parse(regex string) *tokenCtx {
	ctx := &tokenCtx{
		pos:    0,
		tokens: []token{},
	}
	for ctx.pos < len(regex) {
		process(regex, ctx)
		ctx.pos += 1
	}
	return ctx
}

func main() {
	//txt := "a|b[a-z]*"
	txt := "(a|b)[a-z]+"
	fmt.Printf("Input is %v\n", txt)
	ctx := parse(txt)

	PrintTokens(ctx.tokens)
}
