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

func (rp *RegexParser) DoesMatch(matcherString string) bool {
	return rp.engine.DoesMatch(matcherString)
}
