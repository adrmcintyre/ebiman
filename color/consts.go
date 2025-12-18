package color

type Palette byte

const (
	PAL_BLACK        Palette = 0  // blk blk blk blk
	PAL_BLINKY       Palette = 1  // blk wht blu red
	PAL_PINKY        Palette = 3  // blk wht blu pnk
	PAL_INKY         Palette = 5  // blk wht blu cyn
	PAL_CLYDE        Palette = 7  // blk wht blu org
	PAL_PACMAN       Palette = 9  // blk blu red yel
	PAL_PACMAN2      Palette = 9  // ^^
	PAL_GHOST_BONUS  Palette = 9  // ^^ 200,400,800,1600
	PAL_GALAXIAN     Palette = 9  // ^^
	PAL_PACMAN_GREEN Palette = 10 // blk blu red grn
	PAL_14           Palette = 14 // blk wht blk pil
	PAL_SCORE        Palette = 15 // blk red grn wht
	PAL_STRAWBERRY   Palette = 15 // ^^
	PAL_MAZE         Palette = 16 // blk pil blk blu
	PAL_SCARED       Palette = 17 // blk grn blu pil
	PAL_SCARED_FLASH Palette = 18 // blk grn wht red
	PAL_APPLE        Palette = 20 // blk red brn wht
	PAL_CHERRY       Palette = 20 // ^^
	PAL_ORANGE       Palette = 21 // blk org grn brn
	PAL_BELL         Palette = 22 // blk yel stl wht
	PAL_KEY          Palette = 22 // ^^
	PAL_PINEAPPLE    Palette = 23 // blk stl grn wht
	PAL_GATE         Palette = 24 // blk cyn pnk yel
	PAL_EYES         Palette = 25 // blk wht blu blk
	PAL_26           Palette = 26 // blk pil blk blu - maze special?
	PAL_TUNNEL       Palette = 27 // blk pil blk blu
	PAL_29           Palette = 29 // blk wht pil red
	PAL_30           Palette = 30 // blk wht blu pil
	PAL_MAZE_FLASH   Palette = 31 // blk pil blk wht
)
