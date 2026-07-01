package parser

import (
	"slices"
)

func (e *Engine) PrintNFA() {
	e.nfa.Print()
}

func (e *Engine) epsilonClosure(curStates []*state) {
	e.queue = e.queue[:0]
	e.workingSet = e.workingSet[:0]
	clear(e.seen)

	for _, c := range curStates {
		if ok := e.seen[c.id]; !ok {
			e.workingSet = append(e.workingSet, c)
			e.queue = append(e.queue, c)
		}
		e.seen[c.id] = true
	}

	clear(e.seen)

	for len(e.queue) != 0 {
		poppedState := e.queue[0]
		e.queue[0] = nil
		e.queue = e.queue[1:]

		if ok := e.seen[poppedState.id]; ok {
			continue
		}

		e.seen[poppedState.id] = true

		for _, ed := range poppedState.transitions {
			if ed.edge.kind == edgeEpsilon {
				e.workingSet = append(e.workingSet, ed.state)
				e.queue = append(e.queue, ed.state)
			}
		}
	}
}

func (e *Engine) MatchString(matcher string) bool {
	e.epsilonClosure([]*state{e.nfa.start})
	e.next = e.next[:0]

	for i := 0; i < len(matcher); i++ {
		next := []*state{}
		ch := matcher[i]

		for _, s := range e.workingSet {
			for _, ed := range s.transitions {
				switch ed.edge.kind {
				case edgeByte:
					if ed.edge.bitMap[ch/64]&(1<<(ch%64)) != 0 {
						next = append(next, ed.state)
					}
				case edgeLiteral:
					if ed.edge.ch == ch {
						next = append(next, ed.state)
					}
				case edgeEpsilon:
				default:
					panic("Should never happen")
				}
			}
		}

		e.epsilonClosure(next)
	}

	return slices.Contains(e.workingSet, e.nfa.end)
}

func (e *Engine) FindPrefixMatch(start int, matcher string) (string, bool) {
	e.epsilonClosure([]*state{e.nfa.start})
	e.next = e.next[:0]

	for end := start; end < len(matcher); end++ {
		ch := matcher[end]

		for _, s := range e.workingSet {
			for _, ed := range s.transitions {
				if ed.edge.kind == edgeLiteral && ed.edge.ch == ch {
					e.next = append(e.next, ed.state)
				}
			}
		}

		e.epsilonClosure(e.next)
		if len(e.workingSet) == 0 {
			break
		}

		if slices.Contains(e.workingSet, e.nfa.end) {
			return matcher[start : end+1], true
		}
	}
	return "", false
}

func (e *Engine) FindFirstMatch(matcher string) (string, bool) {
	for start := 0; start < len(matcher); start++ {
		m, ok := e.FindPrefixMatch(start, matcher)
		if ok {
			return m, ok
		}
	}

	return "", false
}
