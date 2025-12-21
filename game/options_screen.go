package game

import (
	"fmt"

	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/input"
	"github.com/adrmcintyre/poweraid/option"
	"github.com/adrmcintyre/poweraid/tile"
)

// AnimOptionsScreen is an animator coroutine for the game's menu / start screen.
func (g *Game) AnimOptionsScreen(frame int) (nextFrame int, delay int) {
	next := frame + 1

	const menuLeft = 2
	const menuTop = 8
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
			"PLAYERS    - ",
			[]Option{
				{"1 PLAYER *", option.GAME_MODE_1P},
				{"2 PLAYER  ", option.GAME_MODE_2P},
			},
			0,
			&g.Options.GameMode,
		},
		{
			"LIVES      - ",
			[]Option{
				{"3 *", 3},
				{"5  ", 5},
				{"10 ", 10},
			},
			0,
			&g.Options.Lives,
		},
		{
			"DIFFICULTY - ",
			[]Option{
				{"EASY    ", option.DIFFICULTY_EASY},
				{"NORMAL *", option.DIFFICULTY_MEDIUM},
				{"HARD    ", option.DIFFICULTY_HARD},
			},
			0,
			&g.Options.Difficulty,
		},
		/*
			{
				"NUM GHOSTS - ",
				[]Option{
					{"1  ", 1},
					{"2  ", 2},
					{"3  ", 3},
					{"4 *", 4},
				},
				0,
				&g.Options.MaxGhosts,
			},
		*/
		{
			"GHOST AI   - ",
			[]Option{
				{"OFF ", option.GHOST_AI_OFF},
				{"ON *", option.GHOST_AI_ON},
			},
			0,
			&g.Options.GhostAi,
		},
		/*
			{
				"FRAME RATE - ",
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

	v := &g.Video

	switch frame {
	case 0:
		g.LevelState.DemoMode = false
		g.HideActors()

		v.ClearTiles() // zero out splash screen cruft

		// TODO disable flashing score

		v.SetCursor(6, menuTop-3)
		v.WriteString("OPTIONS MENU", color.PAL_SCORE)

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

		v.SetCursor(6, 19)
		v.WriteString("1 OR 2 PLAYERS", color.PAL_INKY)

		v.SetCursor(2, 22)
		msg := fmt.Sprintf("BONUS LIFE AT %d ", data.EXTRA_LIFE_SCORE)

		v.WriteString(msg, color.PAL_30)
		v.WriteTiles([]tile.Tile{tile.PTS, tile.PTS + 1, tile.PTS + 2}, color.PAL_30)

		v.SetCursor(4, 25)
		v.WriteString("ARROW KEYS TO MOVE", color.PAL_SCORE)

		v.SetCursor(4, 29)
		v.WriteString("  SPACE TO START  ", color.PAL_BLINKY)

		g.StartMenuIndex = 0

		return next, 0

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
				return 0, 0
			}
			g.StartMenuIndex = menuIndex
		}

		return frame, 0

	default:
		return 0, 0
	}
}
