package lexer

import (
	"fmt"
)

func ProduceTokens(regex string) (error, []Token) {
	ctx := &TokenCtx{
		Pos:    0,
		Tokens: []Token{},
	}

	for ctx.Pos < len(regex) {
		err := process(regex, ctx)
		if err != nil {
			return err, nil
		}
		ctx.Pos += 1
	}

	return nil, ctx.Tokens
}

func process(regex string, ctx *TokenCtx) error {
	cur := regex[ctx.Pos]
	switch cur {
	case '[':
		err := processBracket(ctx, regex)
		if err != nil {
			return err
		}
	case '+', '*', '?':
		err := processRepeat(ctx, regex)
		if err != nil {
			return err
		}
	case '(':
		groupCtx := &TokenCtx{
			Pos:    ctx.Pos,
			Tokens: []Token{},
		}
		err := processGroup(groupCtx, regex)
		if err != nil {
			return err
		}
		ctx.Tokens = append(ctx.Tokens, Token{
			TokenType: Group,
			Value:     groupCtx.Tokens,
		})
		ctx.Pos = groupCtx.Pos
	case '{':
		err := processRepeatLimits(ctx, regex)
		if err != nil {
			return err
		}
	case '|':
		err := processOr(ctx, regex)
		if err != nil {
			return err
		}
	default:
		ctx.Tokens = append(ctx.Tokens, Token{
			TokenType: Literal,
			Value:     cur,
		})
	}
	return nil
}

func processOr(ctx *TokenCtx, regex string) error {
	rhsContext := &TokenCtx{
		Pos:    ctx.Pos + 1,
		Tokens: []Token{},
	}
	for rhsContext.Pos < len(regex) && regex[rhsContext.Pos] != ')' {
		err := process(regex, rhsContext)
		if err != nil {
			return err
		}
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
	return nil
}

func processRepeatLimits(ctx *TokenCtx, regex string) error {
	if len(ctx.Tokens) == 0 {
		return fmt.Errorf("Repeat has no target")
	}
	lower, upper, err := getBracketLimits(ctx, regex)
	if err != nil {
		return err
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
	return nil
}

func processGroup(groupCtx *TokenCtx, regex string) error {
	groupCtx.Pos += 1
	for regex[groupCtx.Pos] != ')' {
		err := process(regex, groupCtx)
		if err != nil {
			return err
		}
		groupCtx.Pos += 1
	}
	return nil
}

func processRepeat(ctx *TokenCtx, regex string) error {
	if len(ctx.Tokens) == 0 {
		panic("Repeat has no target")
	}
	ch := regex[ctx.Pos]
	lower, upper, err := getLimits(ch)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
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
	return nil
}

func processBracket(ctx *TokenCtx, regex string) error {
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

	return nil
}
