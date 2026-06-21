package parser

import (
	"fmt"
	"github.com/exgene/rp/internal/lexer"
)

type state struct {
	// Representation of one to many current state => possible transitions.
	id          uint16
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
	case edgeEpsilon:
		return ("{Epsillon}")
	case edgeLiteral:
		return ("{Edge}")
	default:
		panic(fmt.Sprintf("unexpected main.edgeType: %#v", e))
	}
}

const (
	edgeLiteral edgeType = iota
	edgeEpsilon edgeType = iota
)

type edge struct {
	kind edgeType
	ch   byte
}

type nfa struct {
	// so instead of having types, you create the NFA and store start
	// and end states directly into the NFA structure, maybe not the heavy states
	// but a pointer to it is also fine.
	start *state
	end   *state
}

type compiler struct {
	nextID uint16
}

func (c *compiler) next() uint16 {
	c.nextID += 1
	return c.nextID - 1
}

type Engine struct {
	nfa *nfa
}

func Build(regex string) (error, *Engine) {
	err, toks := lexer.ProduceTokens(regex)
	if err != nil {
		return err, nil
	}
	c := &compiler{}
	err, nfa := c.buildNFA(toks)
	if err != nil {
		return err, nil
	}
	return nil, &Engine{
		nfa: nfa,
	}
}

func (n *nfa) build(c *compiler) {
	start := state{
		id:          c.next(),
		transitions: []*transition{},
	}
	end := state{
		id:          c.next(),
		transitions: []*transition{},
	}
	n.start = &start
	n.end = &end
}

// should link the left NFA with the RIGHT one using
// epsillon connection between the two
func (left *nfa) link(right *nfa) {
	left.end.transitions = append(left.end.transitions, &transition{
		edge:  edge{kind: edgeEpsilon},
		state: right.start,
	})
	left.end = right.end
}

func (c *compiler) buildNFA(toks []lexer.Token) (error, *nfa) {
	var prevNFA *nfa = nil
	var returnPTR *nfa = nil

	for _, tok := range toks {
		if prevNFA == nil {
			prevNFA = &nfa{}
			prevNFA.build(c)
			err := c.handleTok(&tok, prevNFA)
			if err != nil {
				return err, nil
			}
			returnPTR = prevNFA
		} else {
			newNFA := nfa{}
			newNFA.build(c)
			err := c.handleTok(&tok, &newNFA)
			if err != nil {
				return err, nil
			}
			prevNFA.link(&newNFA)
			prevNFA = &newNFA
			returnPTR.end = newNFA.end
		}
	}

	return nil, returnPTR
}

func (s *state) Print() {
	fmt.Printf("state:%d\n", s.id)
}

func (n nfa) Print() {
	// map / set to prevent recursive infinite printing
	id := map[*state]bool{n.start: true}
	queue := []*state{n.start}
	fmt.Print("NFA main obj start and end =>\n")
	n.start.Print()
	n.end.Print()

	for len(queue) > 0 {
		cur := queue[0]
		queue[0] = nil
		queue = queue[1:]

		label := ""
		if cur == n.start {
			label = " (start) "
		}
		if cur == n.end {
			label = " (end) "
		}
		fmt.Printf("State %d%s:\n", cur.id, label)

		for _, t := range cur.transitions {
			if _, seen := id[t.state]; !seen {
				id[t.state] = true
				queue = append(queue, t.state)
			}
			switch t.edge.kind {
			case edgeEpsilon:
				fmt.Printf("  --%s--> State %d\n", t.edge.kind.String(), t.state.id)
			case edgeLiteral:
				fmt.Printf("  --%c--> State %d\n", t.edge.ch, t.state.id)
			default:
				panic(fmt.Sprintf("unexpected main.edgeType: %#v", t.edge.kind))
			}
		}

		if len(cur.transitions) == 0 {
			fmt.Println("  (no outgoing transitions)")
		}
	}
}

func (c *compiler) handleTok(tok *lexer.Token, n *nfa) error {
	// build NFA for the given token
	switch tok.TokenType {
	case lexer.Bracket:
		// you don't need ok because map only iterates over truthy values...
		// for ch, ok := range Not required
		for ch := range tok.Value.(map[uint8]bool) {
			n.start.transitions = append(n.start.transitions, &transition{
				edge:  edge{kind: edgeLiteral, ch: byte(ch)},
				state: n.end,
			})
		}
	case lexer.Group, lexer.GroupUncaptured:
		err, inner := c.buildNFA(tok.Value.([]lexer.Token))
		if err != nil {
			return err
		}
		n.start.transitions = append(n.start.transitions, &transition{
			edge:  edge{kind: edgeEpsilon},
			state: inner.start,
		})
		inner.end.transitions = append(inner.end.transitions, &transition{
			edge:  edge{kind: edgeEpsilon},
			state: n.end,
		})
	case lexer.Literal:
		n.start.transitions = append(n.start.transitions, &transition{
			edge: edge{
				kind: edgeLiteral,
				ch:   tok.Value.(byte),
			},
			state: n.end,
		})
	case lexer.Or:
		left := tok.Value.([]lexer.Token)[0]
		right := tok.Value.([]lexer.Token)[1]
		err, leftNFA := c.buildNFA([]lexer.Token{left})
		if err != nil {
			return err
		}
		err, rightNFA := c.buildNFA([]lexer.Token{right})
		if err != nil {
			return err
		}
		n.start.transitions = append(n.start.transitions, &transition{
			edge:  edge{kind: edgeEpsilon},
			state: leftNFA.start,
		})
		n.start.transitions = append(n.start.transitions, &transition{
			edge:  edge{kind: edgeEpsilon},
			state: rightNFA.start,
		})
		leftNFA.end.transitions = append(leftNFA.end.transitions, &transition{
			edge:  edge{kind: edgeEpsilon},
			state: n.end,
		})
		rightNFA.end.transitions = append(rightNFA.end.transitions, &transition{
			edge:  edge{kind: edgeEpsilon},
			state: n.end,
		})
	case lexer.Repeat:
		payload := tok.Value.(lexer.RepeatPayload)
		if payload.IsStar() {
			err, inner := c.buildNFA([]lexer.Token{payload.Token})
			if err != nil {
				return err
			}
			populateStarWithNFA(inner, n)
		} else if payload.IsPlus() {
			err, inner := c.buildNFA([]lexer.Token{payload.Token})
			if err != nil {
				return err
			}
			err = c.populatePlusWithNFA(inner, n)
			if err != nil {
				return err
			}
		} else {
			err := c.populateCurlyRepeatWithNFA(n, &payload)
			if err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("unexpected main.tokenType: %v", tok.TokenType)
	}
	return nil
}

func (c *compiler) populateCurlyRepeatWithNFA(n *nfa, payload *lexer.RepeatPayload) error {
	var prev *nfa = nil
	for range payload.Min {
		err, inner := c.buildNFA([]lexer.Token{payload.Token})
		if err != nil {
			return err
		}
		if prev == nil {
			n.start.transitions = append(n.start.transitions, &transition{
				edge:  edge{kind: edgeEpsilon},
				state: inner.start,
			})
		} else {
			prev.end.transitions = append(prev.end.transitions, &transition{
				edge:  edge{kind: edgeEpsilon},
				state: inner.start,
			})
		}
		prev = inner
	}

	for i := payload.Min; i < payload.Max; i++ {
		err, inner := c.buildNFA([]lexer.Token{payload.Token})
		if err != nil {
			return err
		}
		if prev == nil {
			n.start.transitions = append(n.start.transitions, &transition{
				edge:  edge{kind: edgeEpsilon},
				state: n.end,
			})
			n.start.transitions = append(n.start.transitions, &transition{
				edge:  edge{kind: edgeEpsilon},
				state: inner.start,
			})
		} else {
			prev.end.transitions = append(prev.end.transitions, &transition{
				edge:  edge{kind: edgeEpsilon},
				state: inner.start,
			})
			prev.end.transitions = append(prev.end.transitions, &transition{
				edge:  edge{kind: edgeEpsilon},
				state: n.end,
			})
		}
		prev = inner
	}
	if prev != nil {
		prev.end.transitions = append(prev.end.transitions, &transition{
			edge:  edge{kind: edgeEpsilon},
			state: n.end,
		})
	}
	return nil
}

func (c *compiler) populatePlusWithNFA(inner *nfa, n *nfa) error {
	n.start.transitions = append(n.start.transitions, &transition{
		edge:  edge{kind: edgeEpsilon},
		state: inner.start,
	})
	inner.end.transitions = append(inner.end.transitions, &transition{
		edge:  edge{kind: edgeEpsilon},
		state: n.end,
	})
	inner.end.transitions = append(inner.end.transitions, &transition{
		edge:  edge{kind: edgeEpsilon},
		state: n.start,
	})
	return nil
}

func populateStarWithNFA(inner *nfa, n *nfa) error {
	n.start.transitions = append(n.start.transitions, &transition{
		edge:  edge{kind: edgeEpsilon},
		state: inner.start,
	})
	n.start.transitions = append(n.start.transitions, &transition{
		edge:  edge{kind: edgeEpsilon},
		state: n.end,
	})
	inner.end.transitions = append(inner.end.transitions, &transition{
		edge:  edge{kind: edgeEpsilon},
		state: n.start,
	})
	return nil
}
