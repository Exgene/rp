package rp

import (
	"testing"

	"github.com/exgene/rp"
)

type value struct {
	input     string
	regex     string
	doesMatch bool
}

func TestE2EBehaviour(t *testing.T) {
	testCases := []struct {
		desc  string
		value []value
	}{
		{
			desc: "Literal Matching",
			value: []value{
				{
					input:     "a",
					regex:     "a",
					doesMatch: true,
				},
				{
					input:     "aaa",
					regex:     "a",
					doesMatch: false,
				},
				{
					input:     "b",
					regex:     "a",
					doesMatch: false,
				},
			},
		},
		{
			desc: "+, *",
			value: []value{
				{
					input:     "a",
					regex:     "a+",
					doesMatch: true,
				},
				{
					input:     "",
					regex:     "a+",
					doesMatch: false,
				},
				{
					input:     "aaaaaaaaaaaaaaaaaa",
					regex:     "a+",
					doesMatch: true,
				},
				{
					input:     "aaaaaaaaaaaaaaaaaa",
					regex:     "a*",
					doesMatch: true,
				},
				{
					input:     "",
					regex:     "a*",
					doesMatch: true,
				},
			},
		},
		{
			desc: "|",
			value: []value{
				{
					input:     "a",
					regex:     "a|b",
					doesMatch: true,
				},
				{
					input:     "b",
					regex:     "a|b",
					doesMatch: true,
				},
				{
					input:     "ab",
					regex:     "a|b",
					doesMatch: false,
				},
			},
		},
		{
			desc: "[]",
			value: []value{
				{
					input:     "a",
					regex:     "[a-z]+",
					doesMatch: true,
				},
				{
					input:     "b",
					regex:     "[a-z]+",
					doesMatch: true,
				},
				{
					input:     "abc",
					regex:     "[a-z]+",
					doesMatch: true,
				},
				{
					input:     "abc1",
					regex:     "[a-z]+",
					doesMatch: false,
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			engine := rp.NewRegexEngine("")
			for _, v := range tC.value {
				engine.Reset(v.regex)
				if engine.DoesMatch(v.input) != v.doesMatch {
					t.Fatalf("Failed -> %s should %v with %s", v.input, v.doesMatch, v.regex)
				}
			}
		})
	}
}
