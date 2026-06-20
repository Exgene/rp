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

func (e *Engine) DoesMatch(matcher string) bool {
	workingSet := epsilonClosure([]*state{e.nfa.start})

	for i := 0; i < len(matcher); i++ {
		next := []*state{}
		ch := matcher[i]

		for _, s := range workingSet {
			for _, e := range s.transitions {
				if e.edge.kind == edgeLiteral && e.edge.val == ch {
					next = append(next, e.state)
				}
			}
		}

		workingSet = epsilonClosure(next)
	}

	return slices.Contains(workingSet, e.nfa.end)
}
