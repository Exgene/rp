package parser

import (
	"fmt"
	"github.com/exgene/rp/internal/lexer"
)

var incr uint16 = 0

func counter() uint16 {
	c := incr
	incr++
	return c
}

type state struct {
	// Representation of one to many current state => possible transitions.
	nu          uint16
	transitions []*transition
}

// Transition object from A => B, It will have the edge value
// as well as the next transition state.
// This allows for A ==> B and A ==> C to be represented.
type transition struct {
	edge  edge
	state *state
}

type edgeType uint8

func (e edgeType) String() string {
	switch e {
	case edgeEpsillon:
		return ("{Epsillon}")
	case edgeLiteral:
		return ("{Edge}")
	default:
		panic(fmt.Sprintf("unexpected main.edgeType: %#v", e))
	}
}

const (
	edgeLiteral  edgeType = iota
	edgeEpsillon edgeType = iota
)

type edge struct {
	edgeType edgeType
	val      any
}

type NFA struct {
	// so instead of having types, you create the NFA and store start
	// and end states directly into the NFA structure, maybe not the heavy states
	// but a pointer to it is also fine.
	start *state
	end   *state
}

type Engine struct {
	nfa *NFA
}

func Build(regex string) *Engine {
	ctx := &lexer.TokenCtx{
		Pos:    0,
		Tokens: []lexer.Token{},
	}

	for ctx.Pos < len(regex) {
		lexer.Process(regex, ctx)
		ctx.Pos += 1
	}

	return &Engine{
		nfa: buildNFA(ctx.Tokens),
	}
}

func (nfa *NFA) Build() {
	start := state{
		nu:          counter(),
		transitions: []*transition{},
	}
	end := state{
		nu:          counter(),
		transitions: []*transition{},
	}
	nfa.start = &start
	nfa.end = &end
}

// should link the left NFA with the RIGHT one using
// epsillon connection between the two
func (left *NFA) Link(right *NFA) {
	left.end.transitions = append(left.end.transitions, &transition{
		edge:  edge{edgeType: edgeEpsillon},
		state: right.start,
	})
	left.end = right.end
}

func buildNFA(toks []lexer.Token) *NFA {
	var prevNFA *NFA = nil
	var returnPTR *NFA = nil

	for _, tok := range toks {
		if prevNFA == nil {
			prevNFA = &NFA{}
			prevNFA.Build()
			handleTok(&tok, prevNFA)
			returnPTR = prevNFA
		} else {
			newNFA := NFA{}
			newNFA.Build()
			handleTok(&tok, &newNFA)
			prevNFA.Link(&newNFA)
			prevNFA = &newNFA
			returnPTR.end = newNFA.end
		}
	}

	return returnPTR
}

func (s *state) Print() {
	fmt.Printf("state:%d\n", s.nu)
}

func (nfa NFA) Print() {
	// map / set to prevent recursive infinite printing
	id := map[*state]bool{nfa.start: true}
	queue := []*state{nfa.start}
	fmt.Print("NFA main obj start and end =>\n")
	nfa.start.Print()
	nfa.end.Print()

	for len(queue) > 0 {
		cur := queue[0]
		queue[0] = nil
		queue = queue[1:]

		label := ""
		if cur == nfa.start {
			label = " (start) "
		}
		if cur == nfa.end {
			label = " (end) "
		}
		fmt.Printf("State %d%s:\n", cur.nu, label)

		for _, t := range cur.transitions {
			if _, seen := id[t.state]; !seen {
				id[t.state] = true
				queue = append(queue, t.state)
			}
			switch t.edge.edgeType {
			case edgeEpsillon:
				fmt.Printf("  --%s--> State %d\n", t.edge.edgeType.String(), t.state.nu)
			case edgeLiteral:
				fmt.Printf("  --%s--> State %d\n", t.edge.val.(string), t.state.nu)
			default:
				panic(fmt.Sprintf("unexpected main.edgeType: %#v", t.edge.edgeType))
			}
		}

		if len(cur.transitions) == 0 {
			fmt.Println("  (no outgoing transitions)")
		}
	}
}

func handleTok(tok *lexer.Token, nfa *NFA) {
	// build NFA for the given token
	switch tok.TokenType {
	case lexer.Bracket:
		// you don't need ok because map only iterates over truthy values...
		// for ch, ok := range Not required
		for ch := range tok.Value.(map[uint8]bool) {
			nfa.start.transitions = append(nfa.start.transitions, &transition{
				edge:  edge{edgeType: edgeLiteral, val: string(ch)},
				state: nfa.end,
			})
		}
	case lexer.Group, lexer.GroupUncaptured:
		inner := buildNFA(tok.Value.([]lexer.Token))
		nfa.start.transitions = append(nfa.start.transitions, &transition{
			edge:  edge{edgeType: edgeEpsillon},
			state: inner.start,
		})
		inner.end.transitions = append(inner.end.transitions, &transition{
			edge:  edge{edgeType: edgeEpsillon},
			state: nfa.end,
		})
	case lexer.Literal:
		nfa.start.transitions = append(nfa.start.transitions, &transition{
			edge: edge{
				edgeType: edgeLiteral,
				val:      string(tok.Value.(byte)),
			},
			state: nfa.end,
		})
	case lexer.Or:
		left := tok.Value.([]lexer.Token)[0]
		right := tok.Value.([]lexer.Token)[1]
		leftNFA := buildNFA([]lexer.Token{left})
		rightNFA := buildNFA([]lexer.Token{right})
		nfa.start.transitions = append(nfa.start.transitions, &transition{
			edge:  edge{edgeType: edgeEpsillon},
			state: leftNFA.start,
		})
		nfa.start.transitions = append(nfa.start.transitions, &transition{
			edge:  edge{edgeType: edgeEpsillon},
			state: rightNFA.start,
		})
		leftNFA.end.transitions = append(leftNFA.end.transitions, &transition{
			edge:  edge{edgeType: edgeEpsillon},
			state: nfa.end,
		})
		rightNFA.end.transitions = append(rightNFA.end.transitions, &transition{
			edge:  edge{edgeType: edgeEpsillon},
			state: nfa.end,
		})
	case lexer.Repeat:
		payload := tok.Value.(lexer.RepeatPayload)
		if payload.IsStar() {
			inner := buildNFA([]lexer.Token{payload.Token})
			populateStarWithNFA(inner, nfa)
		} else if payload.IsPlus() {
			inner := buildNFA([]lexer.Token{payload.Token})
			populatePlusWithNFA(inner, nfa)
		} else {
			populateCurlyRepeatWithNFA(nfa, &payload)
		}

	default:
		panic(fmt.Sprintf("unexpected main.tokenType: %#v", tok.TokenType))
	}
}

func populateCurlyRepeatWithNFA(nfa *NFA, payload *lexer.RepeatPayload) {
	var prev *NFA = nil
	for range payload.Min {
		inner := buildNFA([]lexer.Token{payload.Token})
		if prev == nil {
			nfa.start.transitions = append(nfa.start.transitions, &transition{
				edge:  edge{edgeType: edgeEpsillon},
				state: inner.start,
			})
		} else {
			prev.end.transitions = append(prev.end.transitions, &transition{
				edge:  edge{edgeType: edgeEpsillon},
				state: inner.start,
			})
		}
		prev = inner
	}

	for i := payload.Min; i < payload.Max; i++ {
		inner := buildNFA([]lexer.Token{payload.Token})
		if prev == nil {
			nfa.start.transitions = append(nfa.start.transitions, &transition{
				edge:  edge{edgeType: edgeEpsillon},
				state: nfa.end,
			})
			nfa.start.transitions = append(nfa.start.transitions, &transition{
				edge:  edge{edgeType: edgeEpsillon},
				state: inner.start,
			})
		} else {
			prev.end.transitions = append(prev.end.transitions, &transition{
				edge:  edge{edgeType: edgeEpsillon},
				state: inner.start,
			})
			prev.end.transitions = append(prev.end.transitions, &transition{
				edge:  edge{edgeType: edgeEpsillon},
				state: nfa.end,
			})
		}
		prev = inner
	}
	if prev != nil {
		prev.end.transitions = append(prev.end.transitions, &transition{
			edge:  edge{edgeType: edgeEpsillon},
			state: nfa.end,
		})
	}
}

func populatePlusWithNFA(inner *NFA, nfa *NFA) {
	nfa.start.transitions = append(nfa.start.transitions, &transition{
		edge:  edge{edgeType: edgeEpsillon},
		state: inner.start,
	})
	inner.end.transitions = append(inner.end.transitions, &transition{
		edge:  edge{edgeType: edgeEpsillon},
		state: nfa.end,
	})
	inner.end.transitions = append(inner.end.transitions, &transition{
		edge:  edge{edgeType: edgeEpsillon},
		state: nfa.start,
	})
}

func populateStarWithNFA(inner *NFA, nfa *NFA) {
	nfa.start.transitions = append(nfa.start.transitions, &transition{
		edge:  edge{edgeType: edgeEpsillon},
		state: inner.start,
	})
	nfa.start.transitions = append(nfa.start.transitions, &transition{
		edge:  edge{edgeType: edgeEpsillon},
		state: nfa.end,
	})
	inner.end.transitions = append(inner.end.transitions, &transition{
		edge:  edge{edgeType: edgeEpsillon},
		state: nfa.start,
	})
}
