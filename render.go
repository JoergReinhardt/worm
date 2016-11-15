package main

import (
	"github.com/nsf/termbox-go"
	"strings"
	//"time"
	"unicode/utf8"
)

type digits []string

func newDigits() digits {
	var dig []string

	dig = append(dig, "┏━┓")
	dig = append(dig, "┃┃┃")
	dig = append(dig, "┗━┛")
	dig = append(dig, "╺┓ ")
	dig = append(dig, " ┃ ")
	dig = append(dig, "╺┻╸")
	dig = append(dig, "┏━┓")
	dig = append(dig, "┏━┛")
	dig = append(dig, "┗━╸")
	dig = append(dig, "┏━┓")
	dig = append(dig, "╺━┫")
	dig = append(dig, "┗━┛")
	dig = append(dig, "╻ ╻")
	dig = append(dig, "┗━┫")
	dig = append(dig, "  ╹")
	dig = append(dig, "┏━╸")
	dig = append(dig, "┗━┓")
	dig = append(dig, "┗━┛")
	dig = append(dig, "┏━┓")
	dig = append(dig, "┣━┓")
	dig = append(dig, "┗━┛")
	dig = append(dig, "┏━┓")
	dig = append(dig, "  ┃")
	dig = append(dig, "  ╹")
	dig = append(dig, "┏━┓")
	dig = append(dig, "┣━┫")
	dig = append(dig, "┗━┛")
	dig = append(dig, "┏━┓")
	dig = append(dig, "┗━┫")
	dig = append(dig, "┗━┛")

	return dig
}

type sector struct {
	pos  int
	str  string
	prev *sector
	next *sector
}
type ring struct {
	*sector
}

func newRing() *ring {

	d := newDigits()
	r := &ring{&sector{}}
	fs := &sector{}

	for i := 0; i < 30; i++ {
		i := i
		// allocate new sector
		s := &sector{
			pos: i,
			str: d[i],
		}
		if i != 0 {
			// set new sector as next sector for rings previous sector
			(*r.sector.prev).next = s
			// set rings previous sector as new sectors previous sector
			(*s).prev = (*r.sector).prev
		} else {
			// set first sector for later referral
			fs = s
		}
		if i != 29 {
			// set new element as rings previous, for next iteration
			(*r).prev = s
		} else {
			(*s).next = fs
		}
	}
	// ring is holding th e last element now…
	// so one ring can be forged, to rule them all (muhaha)
	(*r.sector).next = fs
	return r
}

func counter(i int) {
	// initialize counter with one zero

}

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

	// COUNTER
	counter(g.state.points)

	// CHERRY
	termbox.SetCell((*g.cherry).x, (*g.cherry).y, 'O', RED, BLACK)

	// WORM
	// callback closes over SetCell with proper bg & fg, gets x &y by worms
	// render method.
	var fn = func(x, y int, c rune) {
		termbox.SetCell(x, y, c, GREEN, BLACK)
	}
	// renders through callback
	(*g.worm).render(fn)

	// FLUSH
	termbox.Flush()
}
