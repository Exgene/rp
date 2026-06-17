package main

import (
	"fmt"
)

var incr uint16 = 0

func Next() uint16 {
	c := incr
	incr++
	return c
}

// Transition object from A => B, It will have the edge value
// as well as the next transition state.
// This allows for A ==> B and A ==> C to be represented.
type transition struct {
	edge  string
	state *state
}

type state struct {
	// Representation of one to many current state => possible transitions.
	nu          uint16
	transitions []*transition
}

type NFA struct {
	// so instead of having types, you create the NFA and store start
	// and end states directly into the NFA structure, maybe not the heavy states
	// but a pointer to it is also fine.
	start *state
	end   *state
}

func (nfa *NFA) Build() {
	start := state{
		nu:          Next(),
		transitions: []*transition{},
	}
	end := state{
		nu:          Next(),
		transitions: []*transition{},
	}
	nfa.start = &start
	nfa.end = &end
}

// should link the left NFA with the RIGHT one using
// epsillon connection between the two
func (left *NFA) Link(right *NFA) {
	left.end.transitions = append(left.end.transitions, &transition{
		edge:  "eps",
		state: right.start,
	})
	left.end = right.end
}

func Parse(toks []token) *NFA {
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
			fmt.Printf("  --%s--> State %d\n", t.edge, t.state.nu)
		}

		if len(cur.transitions) == 0 {
			fmt.Println("  (no outgoing transitions)")
		}
	}
}

func handleTok(tok *token, nfa *NFA) {
	// build NFA for the given token
	switch tok.tokenType {
	case bracket:
		// you don't need ok because map only iterates over truthy values...
		// for ch, ok := range Not required
		for ch := range tok.value.(map[uint8]bool) {
			nfa.start.transitions = append(nfa.start.transitions, &transition{
				edge:  string(ch),
				state: nfa.end,
			})
		}
	case group, groupUncaptured:
		inner := Parse(tok.value.([]token))
		nfa.start.transitions = append(nfa.start.transitions, &transition{
			edge:  "eps",
			state: inner.start,
		})
		inner.end.transitions = append(inner.end.transitions, &transition{
			edge:  "eps",
			state: nfa.end,
		})
	case literal:
		nfa.start.transitions = append(nfa.start.transitions, &transition{
			edge:  string(tok.value.(byte)),
			state: nfa.end,
		})
	case or:
		left := tok.value.([]token)[0]
		right := tok.value.([]token)[1]
		leftNFA := Parse([]token{left})
		rightNFA := Parse([]token{right})
		nfa.start.transitions = append(nfa.start.transitions, &transition{
			edge:  "eps",
			state: leftNFA.start,
		})
		nfa.start.transitions = append(nfa.start.transitions, &transition{
			edge:  "eps",
			state: rightNFA.start,
		})
		leftNFA.end.transitions = append(leftNFA.end.transitions, &transition{
			edge:  "eps",
			state: nfa.end,
		})
		rightNFA.end.transitions = append(rightNFA.end.transitions, &transition{
			edge:  "eps",
			state: nfa.end,
		})
	case repeat:
		payload := tok.value.(repeatPayload)
		if payload.isStar() {
			inner := Parse([]token{payload.token})
			populateStarWithNFA(inner, nfa)
		} else if payload.isPlus() {
			inner := Parse([]token{payload.token})
			populatePlusWithNFA(inner, nfa)
		} else {
			populateCurlyRepeatWithNFA(nfa, &payload)
		}

	default:
		panic(fmt.Sprintf("unexpected main.tokenType: %#v", tok.tokenType))
	}
}

func populateCurlyRepeatWithNFA(nfa *NFA, payload *repeatPayload) {
	var prev *NFA = nil
	for range payload.min {
		inner := Parse([]token{payload.token})
		if prev == nil {
			nfa.start.transitions = append(nfa.start.transitions, &transition{
				edge:  "eps",
				state: inner.start,
			})
		} else {
			prev.end.transitions = append(prev.end.transitions, &transition{
				edge:  "eps",
				state: inner.start,
			})
		}
		prev = inner
	}

	for i := payload.min; i < payload.max; i++ {
		inner := Parse([]token{payload.token})
		if prev == nil {
			nfa.start.transitions = append(nfa.start.transitions, &transition{
				edge:  "eps",
				state: nfa.end,
			})
			nfa.start.transitions = append(nfa.start.transitions, &transition{
				edge:  "eps",
				state: inner.start,
			})
		} else {
			prev.end.transitions = append(prev.end.transitions, &transition{
				edge:  "eps",
				state: inner.start,
			})
			prev.end.transitions = append(prev.end.transitions, &transition{
				edge:  "eps",
				state: nfa.end,
			})
		}
		prev = inner
	}
	if prev != nil {
		prev.end.transitions = append(prev.end.transitions, &transition{
			edge:  "eps",
			state: nfa.end,
		})
	}
}

func populatePlusWithNFA(inner *NFA, nfa *NFA) {
	nfa.start.transitions = append(nfa.start.transitions, &transition{
		edge:  "eps",
		state: inner.start,
	})
	inner.end.transitions = append(inner.end.transitions, &transition{
		edge:  "eps",
		state: nfa.end,
	})
	inner.end.transitions = append(inner.end.transitions, &transition{
		edge:  "eps",
		state: nfa.start,
	})
}

func populateStarWithNFA(inner *NFA, nfa *NFA) {
	nfa.start.transitions = append(nfa.start.transitions, &transition{
		edge:  "eps",
		state: inner.start,
	})
	nfa.start.transitions = append(nfa.start.transitions, &transition{
		edge:  "eps",
		state: nfa.end,
	})
	inner.end.transitions = append(inner.end.transitions, &transition{
		edge:  "eps",
		state: nfa.start,
	})
}
