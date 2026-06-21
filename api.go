package rp

import (
	"github.com/exgene/rp/internal/parser"
)

type RegexParser struct {
	regexPattern string
	engine       *parser.Engine
}

func NewRegexEngine(pattern string) (error, RegexParser) {
	err, engine := parser.Build(pattern)
	return err, RegexParser{regexPattern: pattern, engine: engine}
}

func (rp *RegexParser) Reset(pattern string) error {
	var err error
	err, rp.engine = parser.Build(pattern)
	if err != nil {
		return err
	}
	rp.regexPattern = pattern
	return nil
}

func (rp *RegexParser) DoesMatch(matcherString string) bool {
	return rp.engine.DoesMatch(matcherString)
}
