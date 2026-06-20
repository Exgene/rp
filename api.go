package rp

import (
	"github.com/exgene/rp/internal/parser"
)

type RegexParser struct {
	regexPattern string
	engine       *parser.Engine
}

func NewRegexEngine(pattern string) RegexParser {
	engine := parser.Build(pattern)
	return RegexParser{regexPattern: pattern, engine: engine}
}

func (rp *RegexParser) Reset(pattern string) {
	rp.engine = parser.Build(pattern)
	rp.regexPattern = pattern
}

func (rp *RegexParser) DoesMatch(matcherString string) bool {
	return rp.engine.DoesMatch(matcherString)
}
