package main

func (state State) Hash() StateHash {
	const prime = 31
	var hash int

	for _, n := range state {
		hash = hash*prime + n
	}
	return StateHash(hash)
}
