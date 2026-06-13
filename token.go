package main

import (
	"fmt"
)

type tokenType uint8

const (
	group           tokenType = iota
	literal         tokenType = iota
	repeat          tokenType = iota
	bracket         tokenType = iota
	or              tokenType = iota
	groupUncaptured tokenType = iota
)

func (r repeatPayload) isStar() bool {
	return r.max == -1 && r.min == 0
}

func (r repeatPayload) isPlus() bool {
	return r.max == -1 && r.min == 1
}

type token struct {
	tokenType tokenType
	value     any
}

type tokenCtx struct {
	pos    int
	tokens []token
}

type repeatPayload struct {
	min   int
	max   int
	token token
}

func PrintTokens(toks []token) {
	for _, t := range toks {
		t.Print("")
	}
}

func (tok token) Print(indent string) {
	switch tok.tokenType {
	case literal:
		fmt.Printf("%s- literal: %q\n", indent, tok.value.(byte))

	case bracket:
		fmt.Printf("%s- bracket: \n", indent)
		for k := range tok.value.(map[uint8]bool) {
			t := string(k)
			fmt.Printf("%v==", t)
		}

	case group:
		fmt.Printf("%s- group:\n", indent)
		for _, t := range tok.value.([]token) {
			t.Print(indent + "  ")
		}

	case groupUncaptured:
		fmt.Printf("%s- uncaptured group:\n", indent)
		for _, t := range tok.value.([]token) {
			t.Print(indent + "  ")
		}

	case repeat:
		payload := tok.value.(repeatPayload)
		fmt.Printf("%s- repeat: min=%d max=%d\n", indent, payload.min, payload.max)
		fmt.Printf("%s  target:\n", indent)
		payload.token.Print(indent + "    ")

	case or:
		fmt.Printf("%s- or:\n", indent)
		branches := tok.value.([]token)
		for i, branch := range branches {
			fmt.Printf("%s  branch %d:\n", indent, i+1)
			branch.Print(indent + "    ")
		}
	default:
		panic(fmt.Sprintf("unexpected main.tokenType: %#v", tok.tokenType))
	}
}
