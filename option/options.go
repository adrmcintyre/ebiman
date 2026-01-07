package option

// Constants for configuring the overall game play in the options screen.
// All constants are untyped ints to simplify the menu code.

const (
	ModeClassic1P  = iota // single player
	ModeClassic2P         // two players (turn based)
	ModeElectric1P        // single player with charges
	ModeElectric2P        // two players with charges
	modeCount
)

var modeTypeString = [modeCount]string{
	"classic",
	"classic",
	"electric",
	"electric",
}

const (
	GhostAiOff = iota // ghosts hunt at random
	GhostAiOn         // ghosts can hunt actively
)

const (
	DifficultyEasy   = iota // easier than normal
	DifficultyMedium        // "classic" mode
	DifficultyHard          // stupidly hard
	difficultyCount
)

var difficultyString = [difficultyCount]string{
	"easy",
	"medium",
	"hard",
}

// Options describes the selected game play options.
type Options struct {
	Mode       int // Mode* number of players and style
	Difficulty int // Difficulty* how hard is the game
	FrameRate  int // (unused) fps
	MaxGhosts  int // (unused) 1, 2, 3 or 4
	GhostAi    int // GhostAi*
	Lives      int // 3, 5 or 10
}

// DefaultOptions returns a sensible default set of game play options.
func DefaultOptions() Options {
	return Options{
		Mode:       ModeElectric1P,
		Difficulty: DifficultyMedium,
		FrameRate:  60,
		MaxGhosts:  4,
		GhostAi:    GhostAiOn,
		Lives:      5,
	}
}

// NumPlayers returns the number of player selected.
func (o *Options) NumPlayers() int {
	switch o.Mode {
	case ModeClassic1P, ModeElectric1P:
		return 1
	case ModeClassic2P, ModeElectric2P:
		return 2
	default:
		return 1
	}
}

// IsElectric returns true if an "electric" game mode is selected.
func (o *Options) IsElectric() bool {
	switch o.Mode {
	case ModeElectric1P, ModeElectric2P:
		return true
	default:
		return false
	}
}

func (o *Options) LeaderboardName() string {
	const sep = "-"
	return modeTypeString[o.Mode] + sep + difficultyString[o.Difficulty]
}
