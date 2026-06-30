package rp

import (
	"github.com/exgene/rp/internal/parser"
)

type RegexParser struct {
	regexPattern string
	engine       *parser.Engine
}

func NewRegexEngine(pattern string) (RegexParser, error) {
	engine, err := parser.Build(pattern)
	return RegexParser{regexPattern: pattern, engine: engine}, err
}

func (rp *RegexParser) Reset(pattern string) error {
	var err error
	rp.engine, err = parser.Build(pattern)
	if err != nil {
		return err
	}
	rp.regexPattern = pattern
	return nil
}

func (rp *RegexParser) MatchString(matcherString string) bool {
	return rp.engine.MatchString(matcherString)
}

func (rp *RegexParser) PrintNFA() {
	rp.engine.PrintNFA()
}

func (rp *RegexParser) Pattern() string {
	return rp.regexPattern
}

func (rp *RegexParser) FindString(matcherString string) (string, bool) {
	return rp.engine.FindFirstMatch(matcherString)
}

func (rp *RegexParser) FindPrefix(matcherString string) (string, bool) {
	return rp.engine.FindPrefixMatch(0, matcherString)
}
