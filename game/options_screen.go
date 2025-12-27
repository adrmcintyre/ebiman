package game

import (
	"fmt"

	"github.com/adrmcintyre/ebiman/color"
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/input"
	"github.com/adrmcintyre/ebiman/option"
	"github.com/adrmcintyre/ebiman/tile"
)

// AnimOptionsScreen is an animator coroutine for the game's menu / start screen.
func (g *Game) AnimOptionsScreen(coro *Coro) bool {
	const menuLeft = 3
	const menuTop = 5
	const menuSpacing = 2

	type Option struct {
		label string
		value int
	}

	type Menu struct {
		label   string
		options []Option
		sel     int
		value   *int
	}

	menus := []Menu{
		{
			"MODE    ",
			[]Option{
				{"1P CLASSIC   ", option.MODE_CLASSIC_1P},
				{"2P CLASSIC   ", option.MODE_CLASSIC_2P},
				{"1P ELECTRIC *", option.MODE_ELECTRIC_1P},
				{"2P ELECTRIC  ", option.MODE_ELECTRIC_2P},
			},
			0,
			&g.Options.Mode,
		},
		{
			"LIVES      ",
			[]Option{
				{"3 *", 3},
				{"5  ", 5},
				{"10 ", 10},
			},
			0,
			&g.Options.Lives,
		},
		{
			"DIFFICULTY ",
			[]Option{
				{"EASY    ", option.DIFFICULTY_EASY},
				{"NORMAL *", option.DIFFICULTY_MEDIUM},
				{"HARD    ", option.DIFFICULTY_HARD},
			},
			0,
			&g.Options.Difficulty,
		},
		{
			"NUM GHOSTS ",
			[]Option{
				{"1  ", 1},
				{"2  ", 2},
				{"3  ", 3},
				{"4 *", 4},
			},
			0,
			&g.Options.MaxGhosts,
		},
		{
			"GHOST AI   ",
			[]Option{
				{"OFF ", option.GHOST_AI_OFF},
				{"ON *", option.GHOST_AI_ON},
			},
			0,
			&g.Options.GhostAi,
		},
		/*
			{
				"FRAME RATE ",
				[]Option{
					{"10 FPS  ", 10},
					{"30 FPS  ", 30},
					{"40 FPS  ", 40},
					{"50 FPS  ", 50},
					{"60 FPS *", 60},
				},
				0,
				&g.Options.FrameRate,
			},
		*/
	}

	v := g.Video

	switch coro.Step() {
	case 0:
		g.LevelState.DemoMode = false
		g.HideActors()

		v.ClearTiles() // zero out splash screen cruft

		// TODO disable flashing score

		for i, menu := range menus {
			v.SetCursor(menuLeft, menuTop+i*menuSpacing)
			v.WriteString(menu.label, color.PAL_SCORE)
			for _, opt := range menu.options {
				if opt.value == *menu.value {
					v.WriteString(opt.label, color.PAL_PACMAN)
					break
				}
			}
		}

		v.SetCursor(6, 17)
		v.WriteString("1 OR 2 PLAYERS", color.PAL_INKY)

		v.SetCursor(2, 20)
		msg := fmt.Sprintf("BONUS LIFE AT %d ", data.EXTRA_LIFE_SCORE)

		v.WriteString(msg, color.PAL_30)
		v.WriteTiles([]tile.Tile{tile.PTS, tile.PTS + 1, tile.PTS + 2}, color.PAL_30)

		v.SetCursor(4, 24)
		v.WriteString("ARROW KEYS TO MOVE", color.PAL_SCORE)

		v.SetCursor(3, 26)
		v.WriteString("O P VOLUME", color.PAL_SCORE)
		v.WriteString(" * ", color.PAL_29)
		v.WriteString("Q QUIT", color.PAL_SCORE)

		v.SetCursor(4, 29)
		v.WriteString("  SPACE TO START  ", color.PAL_BLINKY)

		g.StartMenuIndex = 0

		return coro.Next()

	case 1:
		menuIndex := g.StartMenuIndex

		selected := map[int]int{}

		for i, menu := range menus {
			for optIndex, opt := range menu.options {
				if opt.value == *menu.value {
					selected[i] = optIndex
					break
				}
			}
		}

		menu := menus[menuIndex]
		sel := selected[menuIndex]
		v.SetCursor(menuLeft-1, menuTop+menuIndex*menuSpacing)
		v.WriteChar('*', color.PAL_SCORE)
		v.SetCursor(menuLeft+len(menu.label), menuTop+menuIndex*menuSpacing)
		v.WriteString(menu.options[sel].label, color.PAL_PACMAN)
		v.ClearRight()

		inp := input.GetJoystickInput()

		if inp != input.JOY_NONE {
			v.SetCursor(menuLeft-1, menuTop+menuIndex*menuSpacing)
			v.WriteChar(' ', color.PAL_PACMAN)

			switch inp {
			case input.JOY_UP:
				menuIndex = (menuIndex + len(menus) - 1) % len(menus)
			case input.JOY_DOWN:
				menuIndex = (menuIndex + 1) % len(menus)
			case input.JOY_LEFT:
				sel = (sel + len(menu.options) - 1) % len(menu.options)
				*menu.value = menu.options[sel].value
			case input.JOY_RIGHT:
				sel = (sel + 1) % len(menu.options)
				*menu.value = menu.options[sel].value
			case input.JOY_BUTTON:
				return coro.Stop()
			}
			g.StartMenuIndex = menuIndex
		}

		return coro.Wait(1)

	default:
		return coro.Stop()
	}
}
