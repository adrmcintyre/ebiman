package color

// A Palette identifies a colour palette (a 3 entry lookup table)
type Palette byte

// Symbolic names for each palette - note that some
// palettes are re-used in different contexts, and
// we therefore give them aliases for clarity.
const (
	PalBlack       Palette = 0  // blk blk blk blk
	PalBlinky      Palette = 1  // blk wht blu red
	PalPinky       Palette = 3  // blk wht blu pnk
	PalInky        Palette = 5  // blk wht blu cyn
	PalClyde       Palette = 7  // blk wht blu org
	PalPacman      Palette = 9  // blk blu red yel
	PalPacman2     Palette = 9  // ^^
	PalGhostBonus  Palette = 9  // ^^ 200,400,800,1600
	PalGalaxian    Palette = 9  // ^^
	PalPacmanGreen Palette = 10 // blk blu red grn
	Pal14          Palette = 14 // blk wht blk pil
	PalScore       Palette = 15 // blk red grn wht
	PalStrawberry  Palette = 15 // ^^
	PalMaze        Palette = 16 // blk pil blk blu
	PalScared      Palette = 17 // blk grn blu pil
	PalScaredFlash Palette = 18 // blk grn wht red
	PalApple       Palette = 20 // blk red brn wht
	PalCherry      Palette = 20 // ^^
	PalOrange      Palette = 21 // blk org grn brn
	PalBell        Palette = 22 // blk yel stl wht
	PalKey         Palette = 22 // ^^
	PalPineapple   Palette = 23 // blk stl grn wht
	PalGate        Palette = 24 // blk cyn pnk yel
	PalEyes        Palette = 25 // blk wht blu blk
	Pal26          Palette = 26 // blk pil blk blu - maze special?
	PalTunnel      Palette = 27 // blk pil blk blu
	Pal29          Palette = 29 // blk wht pil red
	Pal30          Palette = 30 // blk wht blu pil
	PalMazeFlash   Palette = 31 // blk pil blk wht
)
