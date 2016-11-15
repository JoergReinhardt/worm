package main

import (
	"strings"
	"unicode/utf8"

	"github.com/nsf/termbox-go"
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
	msg = append(msg, "                                       ")
	msg = append(msg, "      Feel free to resize screen       ")
	msg = append(msg, "      while playing…                   ")

	//var msg = []string{"test"}
	// get widh and hight of current board
	w, h := termbox.Size()

	// calculate size of boarders around message
	tb := (h - len(msg)) / 2 // top boarder
	lb := (w - 39) / 2       // left boarder
	// painted, painted, painted… painted black

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

var points = newDigit()

func counter() {
	var msg = points.stringRows()

	//var msg = []string{"test"}
	// get widh and hight of current board
	w, _ := termbox.Size()

	// calculate size of boarders around message
	tb := 2                   // top boarder
	lb := w - len(msg[0]) + 2 // left boarder
	// painted, painted, painted… painted black

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

// rendering happens in animation cycle intervalls and gets called by run, once
// per cycle
func render(g *game) {
	// painted, painted, painted… painted black
	termbox.Clear(0, 0)

	// COUNTER
	counter()

	// CHERRY
	termbox.SetCell((*g.cherry).x, (*g.cherry).y, 'O', RED, 0)

	// WORM
	// callback closes over SetCell with proper bg & fg, gets x &y by worms
	// render method.
	var fn = func(x, y int, c rune) {
		termbox.SetCell(x, y, c, GREEN, 0)
	}
	// renders through callback
	(*g.worm).render(fn)
	// FLUSH
	termbox.Flush()

}
