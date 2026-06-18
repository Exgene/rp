package lexer

import (
	"fmt"
)

func Process(regex string, ctx *TokenCtx) {
	cur := regex[ctx.Pos]
	switch cur {
	case '[':
		processBracket(ctx, regex)
	case '+', '*', '?':
		processRepeat(ctx, regex)
	case '(':
		groupCtx := &TokenCtx{
			Pos:    ctx.Pos,
			Tokens: []Token{},
		}
		processGroup(groupCtx, regex)
		ctx.Tokens = append(ctx.Tokens, Token{
			TokenType: Group,
			Value:     groupCtx.Tokens,
		})
		ctx.Pos = groupCtx.Pos
	case '{':
		processRepeatLimits(ctx, regex)
	case '|':
		processOr(ctx, regex)
	default:
		ctx.Tokens = append(ctx.Tokens, Token{
			TokenType: Literal,
			Value:     cur,
		})
	}
}

func processOr(ctx *TokenCtx, regex string) {
	rhsContext := &TokenCtx{
		Pos:    ctx.Pos + 1,
		Tokens: []Token{},
	}
	for rhsContext.Pos < len(regex) && regex[rhsContext.Pos] != ')' {
		Process(regex, rhsContext)
		rhsContext.Pos++
	}

	left := Token{
		TokenType: GroupUncaptured,
		Value:     ctx.Tokens,
	}
	right := Token{
		TokenType: GroupUncaptured,
		Value:     rhsContext.Tokens,
	}

	ctx.Pos = rhsContext.Pos - 1
	ctx.Tokens = []Token{{
		TokenType: Or,
		Value:     []Token{left, right},
	}}
}

func processRepeatLimits(ctx *TokenCtx, regex string) {
	if len(ctx.Tokens) == 0 {
		panic("Repeat has no target")
	}
	lower, upper := getBracketLimits(ctx, regex)
	lastToken := ctx.Tokens[len(ctx.Tokens)-1]
	ctx.Tokens[len(ctx.Tokens)-1] = Token{
		TokenType: Repeat,
		Value: RepeatPayload{
			Min:   lower,
			Max:   upper,
			Token: lastToken,
		},
	}
}

func processGroup(groupCtx *TokenCtx, regex string) {
	groupCtx.Pos += 1
	for regex[groupCtx.Pos] != ')' {
		Process(regex, groupCtx)
		groupCtx.Pos += 1
	}
}

func processRepeat(ctx *TokenCtx, regex string) {
	if len(ctx.Tokens) == 0 {
		panic("Repeat has no target")
	}
	ch := regex[ctx.Pos]
	lower, upper, err := getLimits(ch)
	if err != nil {
		fmt.Printf("Error: %v", err)
		panic("Shouldn't happend: ")
	}
	lastToken := ctx.Tokens[len(ctx.Tokens)-1]
	ctx.Tokens[len(ctx.Tokens)-1] = Token{
		TokenType: Repeat,
		Value: RepeatPayload{
			Min:   lower,
			Max:   upper,
			Token: lastToken,
		},
	}
}

func processBracket(ctx *TokenCtx, regex string) {
	ctx.Pos += 1
	var literals []string
	for ctx.Pos < len(regex) && regex[ctx.Pos] != ']' {
		ch := regex[ctx.Pos]
		if ch == '-' {
			prev := literals[len(literals)-1][0]
			next := regex[ctx.Pos+1]
			literals[len(literals)-1] = fmt.Sprintf("%c%c", prev, next)
		} else {
			literals = append(literals, fmt.Sprintf("%c", ch))
		}
		ctx.Pos++
	}
	// No set in go thats crazy
	literalMap := map[uint8]bool{}
	for _, v := range literals {
		for i := v[0]; i <= v[len(v)-1]; i++ {
			literalMap[i] = true
		}
	}

	ctx.Tokens = append(ctx.Tokens, Token{
		TokenType: Bracket,
		Value:     literalMap,
	})
}
