package main

import (
	"fmt"
)

func process(regex string, ctx *tokenCtx) {
	cur := regex[ctx.pos]
	switch cur {
	case '[':
		processBracket(ctx, regex)
	case '+', '*', '?':
		processRepeat(ctx, regex)
	case '(':
		groupCtx := &tokenCtx{
			pos:    ctx.pos,
			tokens: []token{},
		}
		processGroup(groupCtx, regex)
		ctx.tokens = append(ctx.tokens, token{
			tokenType: group,
			value:     groupCtx.tokens,
		})
		ctx.pos = groupCtx.pos
	case '{':
		processRepeatLimits(ctx, regex)
	case '|':
		processOr(ctx, regex)
	default:
		ctx.tokens = append(ctx.tokens, token{
			tokenType: literal,
			value:     cur,
		})
	}
}

func processOr(ctx *tokenCtx, regex string) {
	rhsContext := &tokenCtx{
		pos:    ctx.pos + 1,
		tokens: []token{},
	}
	for rhsContext.pos < len(regex) && regex[rhsContext.pos] != ')' {
		process(regex, rhsContext)
		rhsContext.pos++
	}

	left := token{
		tokenType: groupUncaptured,
		value:     ctx.tokens,
	}
	right := token{
		tokenType: groupUncaptured,
		value:     rhsContext.tokens,
	}

	ctx.pos = rhsContext.pos - 1
	ctx.tokens = []token{{
		tokenType: or,
		value:     []token{left, right},
	}}
}

func processRepeatLimits(ctx *tokenCtx, regex string) {
	if len(ctx.tokens) == 0 {
		panic("Repeat has no target")
	}
	lower, upper := getBracketLimits(ctx, regex)
	lastToken := ctx.tokens[len(ctx.tokens)-1]
	ctx.tokens[len(ctx.tokens)-1] = token{
		tokenType: repeat,
		value: repeatPayload{
			min:   lower,
			max:   upper,
			token: lastToken,
		},
	}
}

func processGroup(groupCtx *tokenCtx, regex string) {
	groupCtx.pos += 1
	for regex[groupCtx.pos] != ')' {
		process(regex, groupCtx)
		groupCtx.pos += 1
	}
}

func processRepeat(ctx *tokenCtx, regex string) {
	if len(ctx.tokens) == 0 {
		panic("Repeat has no target")
	}
	ch := regex[ctx.pos]
	lower, upper, err := getLimits(ch)
	if err != nil {
		fmt.Printf("Error: %v", err)
		panic("Shouldn't happend: ")
	}
	lastToken := ctx.tokens[len(ctx.tokens)-1]
	ctx.tokens[len(ctx.tokens)-1] = token{
		tokenType: repeat,
		value: repeatPayload{
			min:   lower,
			max:   upper,
			token: lastToken,
		},
	}
}

func processBracket(ctx *tokenCtx, regex string) {
	ctx.pos += 1
	var literals []string
	for ctx.pos < len(regex) && regex[ctx.pos] != ']' {
		ch := regex[ctx.pos]
		if ch == '-' {
			prev := literals[len(literals)-1][0]
			next := regex[ctx.pos+1]
			literals[len(literals)-1] = fmt.Sprintf("%c%c", prev, next)
		} else {
			literals = append(literals, fmt.Sprintf("%c", ch))
		}
		ctx.pos++
	}
	// No set in go thats crazy
	literalMap := map[uint8]bool{}
	for _, v := range literals {
		for i := v[0]; i <= v[len(v)-1]; i++ {
			literalMap[i] = true
		}
	}

	ctx.tokens = append(ctx.tokens, token{
		tokenType: bracket,
		value:     literalMap,
	})
}
