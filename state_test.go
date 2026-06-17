package main

import "testing"

func Test(t *testing.T) {
	testCases := []struct {
		desc string
		toks []token
	}{
		{
			desc: "basic literal test",
			toks: []token{
				{
					tokenType: literal,
					value:     byte('a'),
				},
				{
					tokenType: literal,
					value:     byte('c'),
				},
				{
					tokenType: literal,
					value:     byte('b'),
				},
			},
		},
		{
			desc: "or tests",
			toks: []token{
				{
					tokenType: or,
					value: []token{
						{
							tokenType: literal,
							value:     byte('a'),
						},
						{
							tokenType: literal,
							value:     byte('c'),
						},
					},
				},
			},
		},
		{
			desc: "repeat tokens",
			toks: []token{
				{
					tokenType: repeat,
					value: repeatPayload{
						min: 1,
						max: -1,
						token: token{
							tokenType: literal,
							value:     byte('a'),
						},
					},
				},
			},
		},
		{
			desc: "bracket",
			toks: []token{
				{
					tokenType: bracket,
					value: map[uint8]bool{
						97: true,
						98: true,
						99: true,
					},
				},
			},
		},
		{
			desc: "group",
			toks: []token{
				{
					tokenType: group,
					value: []token{
						{tokenType: literal, value: byte('a')},
						{tokenType: literal, value: byte('b')},
						{tokenType: literal, value: byte('c')},
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// NFA pointer pointing to the first nfa object
			nfa := Parse(tC.toks)
			nfa.Print()
		})
	}
}
