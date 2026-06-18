package parser

import (
	"testing"

	"github.com/exgene/rp/internal/lexer"
)

func BasicParseTestsManualNFAVerification(t *testing.T) {
	testCases := []struct {
		desc string
		toks []lexer.Token
	}{
		{
			desc: "basic literal test",
			toks: []lexer.Token{
				{
					TokenType: lexer.Literal,
					Value:     byte('a'),
				},
				{
					TokenType: lexer.Literal,
					Value:     byte('c'),
				},
				{
					TokenType: lexer.Literal,
					Value:     byte('b'),
				},
			},
		},
		{
			desc: "or tests",
			toks: []lexer.Token{
				{
					TokenType: lexer.Or,
					Value: []lexer.Token{
						{
							TokenType: lexer.Literal,
							Value:     byte('a'),
						},
						{
							TokenType: lexer.Literal,
							Value:     byte('c'),
						},
					},
				},
			},
		},
		{
			desc: "repeat tokens",
			toks: []lexer.Token{
				{
					TokenType: lexer.Repeat,
					Value: lexer.RepeatPayload{
						Min: 1,
						Max: -1,
						Token: lexer.Token{
							TokenType: lexer.Literal,
							Value:     byte('a'),
						},
					},
				},
			},
		},
		{
			desc: "bracket",
			toks: []lexer.Token{
				{
					TokenType: lexer.Bracket,
					Value: map[uint8]bool{
						97: true,
						98: true,
						99: true,
					},
				},
			},
		},
		{
			desc: "group",
			toks: []lexer.Token{
				{
					TokenType: lexer.Group,
					Value: []lexer.Token{
						{TokenType: lexer.Literal, Value: byte('a')},
						{TokenType: lexer.Literal, Value: byte('b')},
						{TokenType: lexer.Literal, Value: byte('c')},
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// NFA pointer pointing to the first nfa object
			nfa := buildNFA(tC.toks)
			nfa.Print()
		})
	}
}
