package rp

import (
	"github.com/exgene/rp/internal/lexer"
	"github.com/exgene/rp/internal/parser"
)

type RegexParser struct {
	regexPattern string
	nfa          *parser.NFA
}

func NewRegexEngine(pattern string) RegexParser {
	nfa := parser.BuildNFA(parse(pattern).Tokens)
	return RegexParser{regexPattern: pattern, nfa: nfa}
}

func (rp *RegexParser) Reset(pattern string) {
	rp.regexPattern = pattern
	rp.nfa = parser.BuildNFA(parse(pattern).Tokens)
}

func (rp *RegexParser) DoesMatch(matcherString string) bool {
	return rp.doesMatch(matcherString)
}

func (rp *RegexParser) doesMatch(matcherString string) bool {
	return parser.DoesMatch(rp.nfa, matcherString)
}

func parse(regex string) *lexer.TokenCtx {
	ctx := &lexer.TokenCtx{
		Pos:    0,
		Tokens: []lexer.Token{},
	}
	for ctx.Pos < len(regex) {
		lexer.Process(regex, ctx)
		ctx.Pos += 1
	}
	return ctx
}
