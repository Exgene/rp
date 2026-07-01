package parser

func (e *Engine) PrintNFA() {
	e.nfa.Print()
}

func (e *Engine) epsilonClosure(curStates []*state) {
	e.queue = e.queue[:0]
	e.workingSet = e.workingSet[:0]
	clear(e.seen)

	for _, c := range curStates {
		if ok := e.seen[c.id]; !ok {
			e.seen[c.id] = true
			e.workingSet = append(e.workingSet, c)
			e.queue = append(e.queue, c)
		}
	}

	for len(e.queue) != 0 {
		poppedState := e.queue[0]
		e.queue[0] = nil
		e.queue = e.queue[1:]

		for _, ed := range poppedState.transitions {
			if ed.edge.kind == edgeEpsilon {
				if ok := e.seen[ed.state.id]; !ok {
					e.seen[ed.state.id] = true
					e.workingSet = append(e.workingSet, ed.state)
					e.queue = append(e.queue, ed.state)
				}
			}
		}
	}
}

func (e *Engine) MatchString(matcher string) bool {
	e.workingSet = e.workingSet[:0]
	e.workingSet = append(e.workingSet, e.startClosure...)

	for i := 0; i < len(matcher); i++ {
		ch := matcher[i]

		next := make([]*state, 0, 4)

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

	for _, s := range e.workingSet {
		if s == e.nfa.end {
			return true
		}
	}
	return false
}

func (e *Engine) FindPrefixMatch(start int, matcher string) (string, bool) {
	e.workingSet = e.workingSet[:0]
	e.workingSet = append(e.workingSet, e.startClosure...)

	for end := start; end < len(matcher); end++ {
		ch := matcher[end]

		next := make([]*state, 0, 4)

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
				}
			}
		}

		e.epsilonClosure(next)
		if len(e.workingSet) == 0 {
			break
		}

		for _, s := range e.workingSet {
			if s == e.nfa.end {
				return matcher[start : end+1], true
			}
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
