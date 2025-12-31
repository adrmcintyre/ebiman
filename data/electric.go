package data

// An ElectricEntry defines the probability of each ghost manipulating
// the charge is it passes over a pill.
type ElectricEntry struct {
	ScaredPct int
	InkyPct   int
	PinkyPct  int
	BlinkyPct int
	ClydePct  int
}

type ElectricConfig struct {
	Easy   ElectricEntry
	Medium ElectricEntry
	Hard   ElectricEntry
}

var Electric = ElectricConfig{
	Easy: ElectricEntry{
		ScaredPct: 30,
		InkyPct:   60,
		PinkyPct:  60,
		BlinkyPct: 40,
		ClydePct:  40,
	},
	Medium: ElectricEntry{
		ScaredPct: 10,
		InkyPct:   70,
		PinkyPct:  70,
		BlinkyPct: 60,
		ClydePct:  60,
	},
	Hard: ElectricEntry{
		ScaredPct: 0,
		InkyPct:   75,
		PinkyPct:  75,
		BlinkyPct: 65,
		ClydePct:  65,
	},
}
