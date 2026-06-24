package parser

import "slices"

func epsilonClosure(curStates []*state) []*state {
	workingSet := []*state{}
	queue := []*state{}
	// ugly go, ugly piece of code fs
	seen := map[uint16]bool{}

	for _, c := range curStates {
		if _, ok := seen[c.id]; !ok {
			workingSet = append(workingSet, c)
			queue = append(queue, c)
		}
		seen[c.id] = true
	}

	clear(seen)

	for len(queue) != 0 {
		poppedState := queue[0]
		queue[0] = nil
		queue = queue[1:]

		if _, ok := seen[poppedState.id]; ok {
			continue
		}

		seen[poppedState.id] = true

		for _, e := range poppedState.transitions {
			if e.edge.kind == edgeEpsilon {
				workingSet = append(workingSet, e.state)
				queue = append(queue, e.state)
			}
		}
	}

	return workingSet
}

func (e *Engine) MatchString(matcher string) bool {
	workingSet := epsilonClosure([]*state{e.nfa.start})

	for i := 0; i < len(matcher); i++ {
		next := []*state{}
		ch := matcher[i]

		for _, s := range workingSet {
			for _, e := range s.transitions {
				if e.edge.kind == edgeLiteral && e.edge.ch == ch {
					next = append(next, e.state)
				}
			}
		}

		workingSet = epsilonClosure(next)
	}

	return slices.Contains(workingSet, e.nfa.end)
}

func (e *Engine) FindPrefixMatch(start int, matcher string) (string, bool) {
	workingSet := epsilonClosure([]*state{e.nfa.start})
	for end := start; end < len(matcher); end++ {
		next := []*state{}
		ch := matcher[end]

		for _, s := range workingSet {
			for _, e := range s.transitions {
				if e.edge.kind == edgeLiteral && e.edge.ch == ch {
					next = append(next, e.state)
				}
			}
		}

		workingSet = epsilonClosure(next)
		if len(workingSet) == 0 {
			break
		}

		if slices.Contains(workingSet, e.nfa.end) {
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
