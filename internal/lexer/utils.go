package lexer

import (
	"strconv"
	"strings"
)

func getLimits(ch byte) (int, int, error) {
	switch ch {
	case '*':
		return 0, -1, nil
	case '+':
		return 1, -1, nil
	case '?':
		return 0, 1, nil

	default:
	}
	return -1, -1, ErrUndefinedRepeatType
}

func getBracketLimits(ctx *TokenCtx, regex string) (int, int, error) {
	start := ctx.Pos + 1
	for ctx.Pos < len(regex) && regex[ctx.Pos] != '}' {
		ctx.Pos++
	}

	if ctx.Pos == len(regex) && regex[ctx.Pos-1] != '}' {
		return 0, 0, ErrMissingClosingFlowerBracket
	}

	boundary := regex[start:ctx.Pos]
	pieces := strings.Split(boundary, ",")
	var min, max int

	if len(pieces) == 1 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			return 0, 0, err
		} else {
			min = value
			max = value
		}
	} else if len(pieces) == 2 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			return 0, 0, err
		} else {
			min = value
		}
		if pieces[1] == "" {
			max = -1
		} else if value, err := strconv.Atoi(pieces[1]); err != nil {
			return 0, 0, err
		} else {
			max = value
		}
	} else {
		return 0, 0, ErrRepeatValueBetweenFlowerBrackets
	}
	return min, max, nil
}
