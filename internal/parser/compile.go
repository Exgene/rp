package parser

import "slices"

func epsilonClosure(curStates []*state) []*state {
	workingSet := []*state{}
	queue := []*state{}
	// ugly go, ugly piece of code fs
	seen := map[uint16]bool{}

	for _, c := range curStates {
		if _, ok := seen[c.nu]; !ok {
			workingSet = append(workingSet, c)
			queue = append(queue, c)
		}
		seen[c.nu] = true
	}

	clear(seen)

	for len(queue) != 0 {
		poppedState := queue[0]
		queue[0] = nil
		queue = queue[1:]

		if _, ok := seen[poppedState.nu]; ok {
			continue
		}

		seen[poppedState.nu] = true

		for _, e := range poppedState.transitions {
			if e.edge.edgeType == edgeEpsillon {
				workingSet = append(workingSet, e.state)
				queue = append(queue, e.state)
			}
		}
	}

	return workingSet
}

func (e *Engine) DoesMatch(matcher string) bool {
	workingSet := epsilonClosure([]*state{e.nfa.start})

	for _, ch := range matcher {
		next := []*state{}

		for _, s := range workingSet {
			for _, e := range s.transitions {
				if e.edge.edgeType == edgeLiteral && e.edge.val.(string) == string(ch) {
					next = append(next, e.state)
				}
			}
		}

		workingSet = epsilonClosure(next)
	}

	return slices.Contains(workingSet, e.nfa.end)
}
