package main

import (
	"github.com/nsf/termbox-go"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	WHITE termbox.Attribute = termbox.ColorWhite
	BLUE  termbox.Attribute = termbox.ColorBlue
	RED   termbox.Attribute = termbox.ColorRed
	GREEN termbox.Attribute = termbox.ColorGreen
	BLACK termbox.Attribute = termbox.ColorBlack
)

type dir uint8

const (
	UP   dir = 0
	DOWN dir = 1 << iota
	LEFT
	RIGHT
)

type GameStat uint8

const (
	INIT GameStat = 0
	RUN  GameStat = 1 << iota
	PAUSE
	GAME_OVER
)

func initScreen() {
	var msg []string
	msg = append(msg, "▄▄      ▄▄                             ")
	msg = append(msg, "██      ██                             ")
	msg = append(msg, "▀█▄ ██ ▄█▀  ▄████▄    ██▄████  ████▄██▄")
	msg = append(msg, " ██ ██ ██  ██▀  ▀██   ██▀      ██ ██ ██")
	msg = append(msg, " ███▀▀███  ██    ██   ██       ██ ██ ██")
	msg = append(msg, " ███  ███  ▀██▄▄██▀   ██       ██ ██ ██")
	msg = append(msg, " ▀▀▀  ▀▀▀    ▀▀▀▀     ▀▀       ▀▀ ▀▀ ▀▀")
	msg = append(msg, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	msg = append(msg, "╻ ╻┏━╸╻  ┏━┓                           ")
	msg = append(msg, "┣━┫┣╸ ┃  ┣━┛╹                          ")
	msg = append(msg, "╹ ╹┗━╸┗━╸╹  ╹                          ")
	msg = append(msg, "                                       ")
	msg = append(msg, "           h, ← : move left            ")
	msg = append(msg, "           l, → : move right           ")
	msg = append(msg, "           j, ↓ : move down            ")
	msg = append(msg, "           k, ↑ : move up              ")
	msg = append(msg, "                                       ")
	msg = append(msg, "              s : start Game           ")
	msg = append(msg, "              p : pause Game           ")
	msg = append(msg, "              q : quit  Game           ")

	//var msg = []string{"test"}
	// get widh and hight of current board
	w, h := termbox.Size()

	// calculate size of boarders around message
	tb := (h - len(msg)) / 2 // top boarder
	lb := (w - 39) / 2       // left boarder
	// painted, painted, painted… painted black

	termbox.Clear(BLACK, BLACK)

	for y, line := range msg {
		y, line := y+tb, line
		strsl := strings.Split(line, "")
		for x, s := range strsl {
			x, s := x+lb, s
			r, _ := utf8.DecodeRuneInString(s)
			termbox.SetCell(x, y, r, WHITE, BLACK)
		}
	}
	termbox.Flush()
}

// rendering happens in animation cycle intervalls and gets called by run, once
// per cycle
func render(g *Game) {
	// painted, painted, painted… painted black
	termbox.Clear(BLACK, BLACK)
	// render cherry
	termbox.SetCell((*g.cherry).x, (*g.cherry).y, 'O', RED, BLACK)

	// callback closes over SetCell with proper bg & fg, gets x &y by worms
	// render method.
	var fn = func(x, y int, c rune) {
		termbox.SetCell(x, y, c, GREEN, BLACK)
	}
	// renders through callback
	(*g.worm).render(fn)
	termbox.Flush()
}

// the gameController runs worm and cherry at the rate required by current worm
// speed
func gameController(g *Game) {
	// set status to run
	// g.state.stat = RUN
	for {
		// INIT SCREEN
		if g.state.stat == INIT {
			initScreen()
			for { // wait for game start or quit
				// check once per render cycle
				time.Sleep(animationSpeed)
				if g.state.stat != INIT {
					break
				}
			}
		}
		// PAUSE MODE
		// if p is pressed, toggle game state and hold loop
		if g.state.stat == PAUSE {
			for {
				// check once per render cycle
				time.Sleep(animationSpeed)
				if g.state.stat != PAUSE {
					break
				}
			}
		}
		// PLAY
		//- detects colissions for next step
		//- ends the game on colission of worm with itself (sets
		//  state.stat to GAME_OVER)
		//- grows the worm and raises its speed on colission
		//  with cherry.
		//- moves the worm one step
		g.play()
		//- accesses game state to read cherry position
		//- use termbox SetCell, to render cherry
		//- pass termbox SetCell to worms render method.
		render(g)
		// wait one worm speed duration cycle until next move (event
		// queue is running parallel meanwhile changing the game state.
		time.Sleep(g.state.speed)
	}
}

func run() { // runs the animation and input event cycles

	//////////////////////////////////////////////////////////////////////////////////////
	/// GAME SETUP
	// Allocate new game, pass termbox size as functional argument
	g := NewGame(termbox.Size)
	// game loop progresses the game in the background
	go gameController(g)

	//////////////////////////////////////////////////////////////////////////////////////
	/// EVENT QUEUE SETUP
	// the event queue yields keyboard events
	eventQueue := make(chan termbox.Event)
	// dplete event queue in the background
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	//////////////////////////////////////////////////////////////////////////////////////
	/// EVENT LOOP
	// event retrieval and render loop. Blocks until game has stopped.
	for {
		ev := <-eventQueue
		if ev.Type == termbox.EventKey {
			switch {
			// set direction for the next move
			case ev.Key == termbox.KeyArrowUp || ev.Ch == 'k':
				g.state.move = UP
			case ev.Key == termbox.KeyArrowDown || ev.Ch == 'j':
				g.state.move = DOWN
			case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'h':
				g.state.move = LEFT
			case ev.Key == termbox.KeyArrowRight || ev.Ch == 'l':
				g.state.move = RIGHT
			case ev.Ch == 's':
				// if on init screen, run, when s is pressed,
				// else ignore
				if g.state.stat == INIT {
					(*g.state).stat = RUN
				}
			case ev.Ch == 'p':
				// toggle game state between pause & run
				if g.state.stat == PAUSE {
					(*g.state).stat = RUN
				} else {
					(*g.state).stat = PAUSE
				}
			case ev.Key == termbox.KeyEsc || ev.Ch == 'q':
				return
			}
		}
		// exit, if game got ended by last move
		if g.state.stat == GAME_OVER {
			// TODO GAME OVER SCREEN
			return
		}

		// render current game state
		render(g)

		// sleep til next animation cycle
		time.Sleep(animationSpeed)
	}
}
