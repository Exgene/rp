package rp

import (
	"github.com/exgene/rp/internal/parser"
)

type RegexParser struct {
	regexPattern string
}

func NewRegexEngine(pattern string) RegexParser {
	parser.Build(pattern)
	return RegexParser{regexPattern: pattern}
}

func (rp *RegexParser) Reset(pattern string) {
	parser.Build(pattern)
	rp.regexPattern = pattern
}

func (rp *RegexParser) DoesMatch(matcherString string) bool {
	return rp.doesMatch(matcherString)
}

func (rp *RegexParser) doesMatch(matcherString string) bool {
	return parser.DoesMatch(matcherString)
}
