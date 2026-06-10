package main

// NextPlayer activates the next player with lives remaining.
func (g *Game) NextPlayer() bool {
	i := g.PlayerNumber
	n := g.Options.NumPlayers()

	for range n {
		i = (i + 1) % n
		if g.Players[i].Lives > 0 {
			g.PlayerNumber = i
			g.Player = &g.Players[i]
			return true
		}
	}
	return false
}
