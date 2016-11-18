package main

import (
	"github.com/nsf/termbox-go"
	"strings"
	"time"
	"unicode/utf8"
)

type view func() []string

func (v view) render(fn func(x, y int, r rune)) {
	msg := v()
	//var msg = []string{"test"}
	// get widh and hight of current board
	w, h := termbox.Size()

	// calculate size of boarders around message
	tb := (h - len(msg)) / 2 // top boarder
	lb := (w - 39) / 2       // left boarder
	// painted, painted, painted… painted black

	//termbox.Clear(0, 0)

	termbox.Clear(0, 0)
	for y, line := range msg {
		y, line := y+tb, line
		strsl := strings.Split(line, "")
		for x, s := range strsl {
			x, s := x+lb, s
			r, _ := utf8.DecodeRuneInString(s)
			termbox.SetCell(x, y, r, WHITE, 0)
		}
	}
	// FLUSH
	termbox.Flush()
}

var initScreen view = func() []string {
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
	msg = append(msg, "                                       ")
	msg = append(msg, "      Feel free to resize screen       ")
	msg = append(msg, "      while playing…                   ")
	return msg
}

var points = newDigit()

// the gameController runs worm and cherry at the rate required by current worm
// speed
func gameController(g *game) {
	// set status to run
	// g.state.stat = RUN
	for {
		// INIT SCREEN
		if g.state.eventState == INIT {
			for { // wait for game start or quit
				var fn = func(x, y int, c rune) {
					termbox.SetCell(x, y, c, WHITE, 0)
				}
				initScreen.render(fn)
				// check once per render cycle
				time.Sleep(animationSpeed)
				if g.state.eventState != INIT {
					break
				}
			}
		}
		// PAUSE MODE
		// if p is pressed, toggle game state and hold loop
		if g.state.eventState == PAUSE {
			for {
				//initScreen()
				// check once per render cycle
				time.Sleep(animationSpeed)
				if g.state.eventState != PAUSE {
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

// rendering happens in animation cycle intervalls and gets called by run, once
// per cycle
func render(g *game) {
	// painted, painted, painted… painted black
	termbox.Clear(0, 0)
	// COUNTER
	fn := func(x, y int, r rune) {
		termbox.SetCell(x, y, r, WHITE, 0)
	}
	(*g.counter).render(fn)
	termbox.Flush()

	// CHERRY
	fn = func(x, y int, r rune) {
		termbox.SetCell(x, y, r, RED, 0)
	}
	(*g.cherry).render(fn)
	termbox.Flush()

	// WORM
	// callback closes over SetCell with proper bg & fg, gets x &y by worms
	// render method.
	fn = func(x, y int, r rune) {
		termbox.SetCell(x, y, r, GREEN, 0)
	}
	// renders through callback
	(*g.worm).render(fn)
	termbox.Flush()

}
