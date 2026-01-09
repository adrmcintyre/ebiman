# ebiman
This is a high quality Pac-Man implementation from first principles written
in Go, using the ebiten gaming library. It tries to strike a reasonable
balance between accuracy and understandability. I have included some additional 
features such as an options screen offering various ways to adjust the game's
difficulty and behaviour. Video & audio post-processing filters aim to give a
flavour of playing on real hardware. If a Nakama server is available, it is 
used for recording high scores.

A novel 'electric' mode is included. Ghosts can positively or negatively charge
the dots in the grid. If the grid's net charge strays too far from neutral, it
enters the danger zone, and eventually becomes fatal. Pac man must restore
balance by eating the charges before this happens.

Tip: while the ghosts are scared after eating a power pill, they cannot lay
additional charges.


## Build & Run for Desktop
```sh
go build .
./ebiman
```

## Build & Run for Browsers with Web Assembly
```sh
./build-wasm.sh
```
Then visit `localhost:8080` in your browser.

## References
* [Understanding Pac Man Ghost Behavior](https://web.archive.org/web/20190903121844/https://gameinternals.com/understanding-pac-man-ghost-behavior)
* [Simon Owen: Pac-Man Emulator](https://web.archive.org/web/20190128054838/https://simonowen.com/articles/pacemu/)
* [Chris Lomont: Pac-Man Emulation Guide](https://www.lomont.org/software/games/pacman/PacmanEmulation.pdf)
* [Jamey Pittman: The Pac-Man Dossier](https://pacman.holenet.info/)
* [Frederic Vecoven: Pacman](http://www.vecoven.com/elec/pacman/pacman.html)
* [Arcade Longplay - Pac-Man (1980) Midway](https://www.youtube.com/watch?v=DBS0yMaSHIo)
* [Pac-Man: Classic Arcade Game Video, History & Game Play Overview](https://web.archive.org/web/20141227025131/https://www.arcadeclassics.net/80s-game-videos/pac-man)
* [Classic Gaming: Pac-Man Technical Information](https://www.classicgaming.cc/classics/pac-man/technical-info)
