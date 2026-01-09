package game

import (
	"fmt"

	"github.com/adrmcintyre/ebiman/color"
	"github.com/adrmcintyre/ebiman/data"
	"github.com/adrmcintyre/ebiman/input"
	"github.com/adrmcintyre/ebiman/message"
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
			"MODE       ",
			[]Option{
				{"1P CLASSIC   ", option.ModeClassic1P},
				{"2P CLASSIC   ", option.ModeClassic2P},
				{"1P ELECTRIC *", option.ModeElectric1P},
				{"2P ELECTRIC  ", option.ModeElectric2P},
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
				{"EASY    ", option.DifficultyEasy},
				{"NORMAL *", option.DifficultyMedium},
				{"HARD    ", option.DifficultyHard},
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
				{"OFF ", option.GhostAiOff},
				{"ON *", option.GhostAiOn},
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
		g.StatusMsg = message.None
		g.PlayerMsg = message.None

		v.ClearTiles() // zero out splash screen cruft

		// TODO disable flashing score

		for i, menu := range menus {
			v.SetCursor(menuLeft, menuTop+i*menuSpacing)
			v.WriteString(menu.label, color.PalScore)
			for _, opt := range menu.options {
				if opt.value == *menu.value {
					v.WriteString(opt.label, color.PalPacman)
					break
				}
			}
		}

		v.SetCursor(6, 17)
		v.WriteString("1 OR 2 PLAYERS", color.PalInky)

		v.SetCursor(2, 20)
		msg := fmt.Sprintf("BONUS LIFE AT %d ", data.ExtraLifeScore)

		v.WriteString(msg, color.Pal30)
		v.WriteTiles([]tile.Tile{tile.Pts, tile.Pts + 1, tile.Pts + 2}, color.Pal30)

		v.SetCursor(6, 24)
		v.WriteString("ARROWS TO MOVE", color.PalScore)

		v.SetCursor(3, 26)
		v.WriteString("O P VOLUME", color.PalScore)
		v.WriteString(" * ", color.Pal29)
		v.WriteString("Q QUIT", color.PalScore)

		v.SetCursor(6, 29)
		v.WriteString("SPACE TO START", color.PalBlinky)

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
		v.WriteChar('*', color.PalScore)
		v.SetCursor(menuLeft+len(menu.label), menuTop+menuIndex*menuSpacing)
		v.WriteString(menu.options[sel].label, color.PalPacman)
		v.ClearRight()

		inp := g.Input.JoystickInput()

		if inp != input.JoyNone {
			v.SetCursor(menuLeft-1, menuTop+menuIndex*menuSpacing)
			v.WriteChar(' ', color.PalPacman)

			switch inp {
			case input.JoyUp:
				menuIndex = (menuIndex + len(menus) - 1) % len(menus)
			case input.JoyDown:
				menuIndex = (menuIndex + 1) % len(menus)
			case input.JoyLeft:
				sel = (sel + len(menu.options) - 1) % len(menu.options)
				*menu.value = menu.options[sel].value
				g.RefreshHighScore()
			case input.JoyRight:
				sel = (sel + 1) % len(menu.options)
				*menu.value = menu.options[sel].value
				g.RefreshHighScore()
			case input.JoyButton:
				return coro.Stop()
			}
			g.StartMenuIndex = menuIndex
		}

		return coro.Wait(1)

	default:
		return coro.Stop()
	}
}
