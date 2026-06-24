package rp

import (
	"testing"

	"github.com/exgene/rp"
)

type boolValue struct {
	input     string
	regex     string
	doesMatch bool
}

type matchValue struct {
	input      string
	regex      string
	isMatching bool
	matching   string
}

func TestE2EMatchPrefix(t *testing.T) {
	testPrefix := []struct {
		desc  string
		value []matchValue
	}{
		{
			desc: "generic prefix match",
			value: []matchValue{
				{
					input:      "",
					regex:      "",
					isMatching: false,
					matching:   "",
				},
			},
		},
	}

	for _, tC := range testPrefix {
		t.Run(tC.desc, func(t *testing.T) {
			engine, err := rp.NewRegexEngine("")
			if err != nil {
				t.Fatalf("Failed to compile regex: %s with error %v", "", err.Error())
			}
			for _, v := range tC.value {
				err := engine.Reset(v.regex)
				if err != nil {
					t.Fatalf("Failed to compile regex: %s with error %v", v.regex, err.Error())
				}
				m, ok := engine.FindString(v.input)
				if ok != v.isMatching || v.matching != m {
					t.Fatalf("Failed => %s doesnt match with %s ::: %v == %v", m, v.matching, ok, v.isMatching)
				}
			}
		})
	}
}

func TestE2EMatchValue(t *testing.T) {
	testMatch := []struct {
		desc  string
		value []matchValue
	}{
		{
			desc: "generic first find",
			value: []matchValue{
				{
					input:      "abab",
					regex:      "ab+",
					matching:   "ab",
					isMatching: true,
				},
				{
					input:      "xxabxx",
					regex:      "ab",
					matching:   "ab",
					isMatching: true,
				},
				{
					input:      "ab",
					regex:      "ab+",
					matching:   "ab",
					isMatching: true,
				},
				{
					input:      "a",
					regex:      "ab+",
					matching:   "",
					isMatching: false,
				},
				{
					input:      "bbb",
					regex:      "ab+",
					matching:   "",
					isMatching: false,
				},
			},
		},
	}

	for _, tC := range testMatch {
		t.Run(tC.desc, func(t *testing.T) {
			engine, err := rp.NewRegexEngine("")
			if err != nil {
				t.Fatalf("Failed to compile regex: %s with error %v", "", err.Error())
			}
			for _, v := range tC.value {
				err := engine.Reset(v.regex)
				if err != nil {
					t.Fatalf("Failed to compile regex: %s with error %v", v.regex, err.Error())
				}
				m, ok := engine.FindString(v.input)
				if ok != v.isMatching || v.matching != m {
					t.Fatalf("Failed => %s doesnt match with %s ::: %v == %v", m, v.matching, ok, v.isMatching)
				}
			}
		})
	}
}

func TestE2EBehaviourDoesMatch(t *testing.T) {
	testBool := []struct {
		desc  string
		value []boolValue
	}{
		{
			desc: "Literal Matching",
			value: []boolValue{
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
			value: []boolValue{
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
			value: []boolValue{
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
			value: []boolValue{
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

	for _, tC := range testBool {
		t.Run(tC.desc, func(t *testing.T) {
			var err error
			engine, err := rp.NewRegexEngine("")
			if err != nil {
				t.Fatalf("Failed to compile regex: %s with error %v", "", err.Error())
			}
			for _, v := range tC.value {
				err := engine.Reset(v.regex)
				if err != nil {
					t.Fatalf("Failed to compile regex: %s with error %v", v.regex, err.Error())
				}
				if engine.MatchString(v.input) != v.doesMatch {
					t.Fatalf("Failed -> %s should %v with %s", v.input, v.doesMatch, v.regex)
				}
			}
		})
	}
}
