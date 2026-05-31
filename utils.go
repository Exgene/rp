package main

import (
	"fmt"
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
	return -1, -1, fmt.Errorf("Err, Undefined var: %v", ch)
}

func getBracketLimits(ctx *tokenCtx, regex string) (int, int) {
	start := ctx.pos + 1
	for regex[ctx.pos] != '}' {
		ctx.pos++
	}

	boundary := regex[start:ctx.pos]
	pieces := strings.Split(boundary, ",")
	var min, max int

	if len(pieces) == 1 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			panic(err.Error())
		} else {
			min = value
			max = value
		}
	} else if len(pieces) == 2 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			panic(err.Error())
		} else {
			min = value
		}
		if pieces[1] == "" {
			max = -1
		} else if value, err := strconv.Atoi(pieces[1]); err != nil {
			panic(err.Error())
		} else {
			max = value
		}
	} else {
		panic("Provide 1, 2 value in between.")
	}
	return min, max
}
