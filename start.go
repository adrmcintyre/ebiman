package main

/*
void start_button_screen()
{
  v.SetCursor(5, 17);
  v.WriteString("PUSH START BUTTON", PAL_CLYDE);

  v.SetCursor(7, 21);
  v.WriteString("1 OR 2 PLAYERS", PAL_INKY);

  v.SetCursor(0, 25);
  char msg[32];
  sprintf(msg, "BONUS PAC MAN FOR %lu ", EXTRA_LIFE_SCORE);
  v.WriteString(msg, PAL_30);
  v.WriteTiles((const byte[]) {
    TILE_PTS, TILE_PTS + 1, TILE_PTS + 2, 0
  }, PAL_30);

  v.SetCursor(0, 29);
  v.WriteTile(TILE_COPYRIGHT, PAL_PINKY);
  v.WriteString(" 2025 MCINTYRE ENTERPRISES", PAL_PINKY);

  struct Option {
    byte value;
    const char* label;
  };

  struct Menu {
    const char* label;
    byte option_count;
    const Option* options;
    byte* value;
  };

  const byte role_server = 0x80;
  byte mode_and_role = opt_game_mode | ((opt_coop_role == COOP_ROLE_SERVER) ? role_server : 0);

  const Option options_mode[] = {
    {GAME_MODE_1P, "1 PLAYER *"},
    {GAME_MODE_2P, "2 PLAYER"},
    {GAME_MODE_COOP | role_server, "CO-OP PLAYER1"},
    {GAME_MODE_COOP, "CO-OP PLAYER2"},
  };
  const Option options_lives[] = {
    {3, "3 *"},
    {5, "5"},
    {10, "10"}
  };
  const Option options_difficulty[] = {
    {0, "EASY"},
    {1, "NORMAL *"},
    {2, "HARD"}
  };
  const Option options_ghosts[] = {
    {1, "1"},
    {2, "2"},
    {3, "3"},
    {4, "4 *"}
  };
  const Option options_ai[] = {
    {0, "OFF"},
    {1, "ON *"}
  };
  const Option options_fps[] = {
    {10, "10 FPS"},
    {30, "30 FPS"},
    {40, "40 FPS"},
    {50, "50 FPS"},
    {60, "60 FPS *"},
  };

  const byte menu_left = 2;
  const byte menu_top = 3;
  const byte menu_spacing = 2;
  const byte menu_count = 6;

  const Menu menus[menu_count] = {
    {"PLAYERS    - ", 4, options_mode,       &mode_and_role  },
    {"LIVES      - ", 3, options_lives,      &opt_lives      },
    {"DIFFICULTY - ", 3, options_difficulty, &opt_difficulty },
    {"NUM GHOSTS - ", 4, options_ghosts,     &opt_max_ghosts },
    {"GHOST AI   - ", 2, options_ai,         &opt_ghost_ai   },
    {"FRAME RATE - ", 5, options_fps,        &opt_frame_rate },
  };

  byte selected[menu_count] = {};

  for(byte i=0; i<menu_count; ++i) {
    const Menu& m = menus[i];
    v.SetCursor(menu_left, menu_top + i * menu_spacing);
    v.WriteString(m.label, PAL_SCORE, true);
    for(byte j=0; j<m.option_count; ++j) {
      if (m.options[j].value == *m.value) {
        selected[i] = j;
        v.WriteString(m.options[j].label, PAL_PACMAN, true);
        break;
      }
    }
  }

  byte menu_index = 0;
  bool start_game = false;

  while(!start_game) {
    const Menu& m = menus[menu_index];
    byte& sel = selected[menu_index];
    v.SetCursor(menu_left-2, menu_top + menu_index * menu_spacing);
    tile_char('*', PAL_SCORE, true);
    v.SetCursor(menu_left + strlen(m.label), menu_top + menu_index * menu_spacing);
    v.WriteString(m.options[sel].label, PAL_PACMAN, true);
    tile_clear_right(true);

    const byte inp = wait_joystick_input();

    v.SetCursor(menu_left-2, menu_top + menu_index * menu_spacing);
    tile_char(' ', PAL_PACMAN, true);

    switch(inp) {
      case JOY_BUTTON:
        start_game = true;
        break;
      case JOY_UP:
        menu_index = (menu_index ? menu_index : menu_count) - 1;
        break;
      case JOY_DOWN:
        if (++menu_index == menu_count) menu_index = 0;
        break;
      case JOY_LEFT:
        sel = (sel ? sel : m.option_count) - 1;
        *m.value = m.options[sel].value;
        break;
      case JOY_RIGHT:
        if (++sel == m.option_count) sel = 0;
        *m.value = m.options[sel].value;
        break;
    }
  }

  opt_game_mode = mode_and_role & ~role_server;
  opt_coop_role = (mode_and_role & role_server) ? COOP_ROLE_SERVER : COOP_ROLE_CLIENT;
}
*/
