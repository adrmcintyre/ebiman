package game

import (
	"fmt"

	"github.com/adrmcintyre/poweraid/color"
	"github.com/adrmcintyre/poweraid/data"
	"github.com/adrmcintyre/poweraid/input"
	"github.com/adrmcintyre/poweraid/tile"
)

func (g *Game) AnimStartButtonScreen(frame int) (nextFrame int, delay int) {
	next := frame + 1

	const menuLeft = 2
	const menuTop = 3
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
				{"1 PLAYER *", GAME_MODE_1P},
				{"2 PLAYER  ", GAME_MODE_2P},
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
				{"EASY    ", DIFFICULTY_EASY},
				{"NORMAL *", DIFFICULTY_NORMAL},
				{"HARD    ", DIFFICULTY_HARD},
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
				{"OFF ", GHOST_AI_OFF},
				{"ON *", GHOST_AI_ON},
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
		v.SetCursor(5, 17)
		v.WriteString("PUSH START BUTTON", color.PAL_CLYDE)

		v.SetCursor(7, 21)
		v.WriteString("1 OR 2 PLAYERS", color.PAL_INKY)

		v.SetCursor(0, 25)
		msg := fmt.Sprintf("BONUS PAC MAN FOR %d ", data.EXTRA_LIFE_SCORE)
		v.WriteString(msg, color.PAL_30)
		v.WriteTiles([]tile.Tile{tile.PTS, tile.PTS + 1, tile.PTS + 2}, color.PAL_30)

		v.SetCursor(0, 29)
		v.WriteTile(tile.COPYRIGHT, color.PAL_PINKY)
		v.WriteString(" 2025 MCINTYRE ENTERPRISES", color.PAL_PINKY)

		for i, m := range menus {
			v.SetCursor(menuLeft, menuTop+i*menuSpacing)
			v.WriteString(m.label, color.PAL_SCORE)
			for _, o := range m.options {
				if o.value == *m.value {
					v.WriteString(o.label, color.PAL_PACMAN)
					break
				}
			}
		}

		g.StartMenuIndex = 0

		return next, 0

	case 1:
		menuIndex := g.StartMenuIndex

		selected := [6]int{}

		for i, m := range menus {
			for j, o := range m.options {
				if o.value == *m.value {
					selected[i] = j
					break
				}
			}
		}

		m := menus[menuIndex]
		sel := selected[menuIndex]
		v.SetCursor(menuLeft-2, menuTop+menuIndex*menuSpacing)
		v.WriteChar('*', color.PAL_SCORE)
		v.SetCursor(menuLeft+len(m.label), menuTop+menuIndex*menuSpacing)
		v.WriteString(m.options[sel].label, color.PAL_PACMAN)
		v.ClearRight()

		inp := input.GetJoystickInput()

		if inp != input.JOY_NONE {
			v.SetCursor(menuLeft-2, menuTop+menuIndex*menuSpacing)
			v.WriteChar(' ', color.PAL_PACMAN)

			switch inp {
			case input.JOY_UP:
				menuIndex = (menuIndex + len(menus) - 1) % len(menus)
			case input.JOY_DOWN:
				menuIndex = (menuIndex + 1) % len(menus)
			case input.JOY_LEFT:
				sel = (sel + len(m.options) - 1) % len(m.options)
				*m.value = m.options[sel].value
			case input.JOY_RIGHT:
				sel = (sel + 1) % len(m.options)
				*m.value = m.options[sel].value
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
