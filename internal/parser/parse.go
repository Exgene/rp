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
		return "{Epsillon}"
	case edgeLiteral:
		return "{Edge}"
	default:
		return "Unknown, This shouldn't happen"
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

func Build(regex string) (*Engine, error) {
	toks, err := lexer.ProduceTokens(regex)
	if err != nil {
		return nil, err
	}
	c := &compiler{}
	nfa, err := c.buildNFA(toks)
	if err != nil {
		return nil, err
	}
	return &Engine{
		nfa: nfa,
	}, nil
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

func (c *compiler) buildNFA(toks []lexer.Token) (*nfa, error) {
	var prevNFA *nfa = nil
	var returnPTR *nfa = nil

	for _, tok := range toks {
		if prevNFA == nil {
			prevNFA = &nfa{}
			prevNFA.build(c)
			err := c.handleTok(&tok, prevNFA)
			if err != nil {
				return nil, err
			}
			returnPTR = prevNFA
		} else {
			newNFA := nfa{}
			newNFA.build(c)
			err := c.handleTok(&tok, &newNFA)
			if err != nil {
				return nil, err
			}
			prevNFA.link(&newNFA)
			prevNFA = &newNFA
			returnPTR.end = newNFA.end
		}
	}

	return returnPTR, nil
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
				// wont be reached ig?
				fmt.Printf("UnknownType --> State %d\n", t.state.id)
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
		m, ok := tok.Value.(map[uint8]bool)
		if !ok {
			return lexer.ErrExpectedValueShapeMismatch
		}
		for ch := range m {
			n.start.transitions = append(n.start.transitions, &transition{
				edge:  edge{kind: edgeLiteral, ch: byte(ch)},
				state: n.end,
			})
		}
	case lexer.Group, lexer.GroupUncaptured:
		tok, ok := tok.Value.([]lexer.Token)
		if !ok {
			return lexer.ErrExpectedValueShapeMismatch
		}
		inner, err := c.buildNFA(tok)
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
		ch, ok := tok.Value.(byte)
		if !ok {
			return lexer.ErrExpectedValueShapeMismatch
		}
		n.start.transitions = append(n.start.transitions, &transition{
			edge: edge{
				kind: edgeLiteral,
				ch:   ch,
			},
			state: n.end,
		})
	case lexer.Or:
		v, ok := tok.Value.([]lexer.Token)
		if !ok {
			return lexer.ErrExpectedValueShapeMismatch
		}
		left := v[0]
		right := v[1]
		leftNFA, err := c.buildNFA([]lexer.Token{left})
		if err != nil {
			return err
		}
		rightNFA, err := c.buildNFA([]lexer.Token{right})
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
		payload, ok := tok.Value.(lexer.RepeatPayload)
		if !ok {
			return lexer.ErrExpectedValueShapeMismatch
		}
		if payload.IsStar() {
			inner, err := c.buildNFA([]lexer.Token{payload.Token})
			if err != nil {
				return err
			}
			populateStarWithNFA(inner, n)
		} else if payload.IsPlus() {
			inner, err := c.buildNFA([]lexer.Token{payload.Token})
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
		return ErrUnexpectedToken
	}
	return nil
}

func (c *compiler) populateCurlyRepeatWithNFA(n *nfa, payload *lexer.RepeatPayload) error {
	var prev *nfa = nil
	for range payload.Min {
		inner, err := c.buildNFA([]lexer.Token{payload.Token})
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

	if payload.IsCurly() && prev != nil {
		prev.end.transitions = append(prev.start.transitions, &transition{
			edge: edge{
				kind: edgeEpsilon,
			},
			state: n.start,
		})

		prev.end.transitions = append(prev.end.transitions, &transition{
			edge:  edge{kind: edgeEpsilon},
			state: n.end,
		})
		return nil
	}

	for i := payload.Min; i < payload.Max; i++ {
		inner, err := c.buildNFA([]lexer.Token{payload.Token})
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
