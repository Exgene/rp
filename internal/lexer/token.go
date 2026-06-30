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

func (r RepeatPayload) IsCurly() bool {
	return r.Max == -1 && r.Min >= 0
}

func (r RepeatPayload) IsPlus() bool {
	return r.Max == -1 && r.Min == 1
}

type Token struct {
	TokenType TokenType
	Value     any
}

type TokenCtx struct {
	Pos    int
	Tokens []Token
}

type RepeatPayload struct {
	Min   int
	Max   int
	Token Token
}

func PrintTokens(toks []Token) error {
	for _, t := range toks {
		err := t.Print("")
		if err != nil {
			return err
		}
	}
	return nil
}

func (tok Token) Print(indent string) error {
	switch tok.TokenType {
	case Literal:
		ch, ok := tok.Value.(byte)
		if !ok {
			return ErrExpectedValueShapeMismatch
		}
		fmt.Printf("%s- literal: %q\n", indent, ch)

	case Bracket:
		fmt.Printf("%s- bracket: \n", indent)
		m, ok := tok.Value.(map[uint8]bool)
		if !ok {
			return ErrExpectedValueShapeMismatch
		}
		for k := range m {
			t := string(k)
			fmt.Printf("%v==", t)
		}

	case Group:
		fmt.Printf("%s- group:\n", indent)
		m, ok := tok.Value.([]Token)
		if !ok {
			return ErrExpectedValueShapeMismatch
		}
		for _, t := range m {
			t.Print(indent + "  ")
		}

	case GroupUncaptured:
		fmt.Printf("%s- uncaptured group:\n", indent)
		m, ok := tok.Value.([]Token)
		if !ok {
			return ErrExpectedValueShapeMismatch
		}
		for _, t := range m {
			t.Print(indent + "  ")
		}

	case Repeat:
		payload, ok := tok.Value.(RepeatPayload)
		if !ok {
			return ErrExpectedValueShapeMismatch
		}
		fmt.Printf("%s- repeat: min=%d max=%d\n", indent, payload.Min, payload.Max)
		fmt.Printf("%s  target:\n", indent)
		payload.Token.Print(indent + "    ")

	case Or:
		fmt.Printf("%s- or:\n", indent)
		branches, ok := tok.Value.([]Token)
		if !ok {
			return ErrExpectedValueShapeMismatch
		}
		for i, branch := range branches {
			fmt.Printf("%s  branch %d:\n", indent, i+1)
			branch.Print(indent + "    ")
		}
	default:
		return ErrUnexpectedToken
	}
	return nil
}
