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
					value:     "a",
				},
				{
					tokenType: literal,
					value:     "c",
				},
				{
					tokenType: literal,
					value:     "b",
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
							value:     "a",
						},
						{
							tokenType: literal,
							value:     "c",
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
							value:     "a",
						},
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
