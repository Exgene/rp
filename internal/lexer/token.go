package lexer

import (
	"fmt"
)

type TokenType uint8

const (
	Group           TokenType = iota
	Literal         TokenType = iota
	Repeat          TokenType = iota
	Bracket         TokenType = iota
	Or              TokenType = iota
	GroupUncaptured TokenType = iota
)

func (r RepeatPayload) IsStar() bool {
	return r.Max == -1 && r.Min == 0
}

func (r RepeatPayload) IsPlus() bool {
	return r.Max == -1 && r.Min == 1
}

type Token struct {
	TokenType TokenType
	Value     any
}

type tokenCtx struct {
	Pos    int
	Tokens []Token
}

type RepeatPayload struct {
	Min   int
	Max   int
	Token Token
}

func PrintTokens(toks []Token) {
	for _, t := range toks {
		t.Print("")
	}
}

func (tok Token) Print(indent string) {
	switch tok.TokenType {
	case Literal:
		fmt.Printf("%s- literal: %q\n", indent, tok.Value.(byte))

	case Bracket:
		fmt.Printf("%s- bracket: \n", indent)
		for k := range tok.Value.(map[uint8]bool) {
			t := string(k)
			fmt.Printf("%v==", t)
		}

	case Group:
		fmt.Printf("%s- group:\n", indent)
		for _, t := range tok.Value.([]Token) {
			t.Print(indent + "  ")
		}

	case GroupUncaptured:
		fmt.Printf("%s- uncaptured group:\n", indent)
		for _, t := range tok.Value.([]Token) {
			t.Print(indent + "  ")
		}

	case Repeat:
		payload := tok.Value.(RepeatPayload)
		fmt.Printf("%s- repeat: min=%d max=%d\n", indent, payload.Min, payload.Max)
		fmt.Printf("%s  target:\n", indent)
		payload.Token.Print(indent + "    ")

	case Or:
		fmt.Printf("%s- or:\n", indent)
		branches := tok.Value.([]Token)
		for i, branch := range branches {
			fmt.Printf("%s  branch %d:\n", indent, i+1)
			branch.Print(indent + "    ")
		}
	default:
		panic(fmt.Sprintf("unexpected main.tokenType: %#v", tok.TokenType))
	}
}
