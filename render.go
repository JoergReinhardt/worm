package main

import (
	"github.com/nsf/termbox-go"
	"strings"
	"time"
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

func (r *ring) stringRow() string { return r.str }
func (r *ring) stringDig() string {
	ret := r.prev.prev.str + "\n"
	ret = ret + r.prev.str + "\n"
	ret = ret + r.str + "\n"
	return ret
}
func (r *ring) nextRow() {
	(*r).sector = (*r).next
}
func (r *ring) prevRow() {
	(*r).sector = (*r).prev
}
func (r *ring) nextDigit() {
	for i := 0; i < 3; i++ {
		(*r).nextRow()
	}
}
func (r *ring) prevDigit() {
	for i := 0; i < 3; i++ {
		(*r).prevRow()
	}
}
func (r *ring) nextDigitAnimated() {
	for i := 0; i < 3; i++ {
		(*r).nextRow()
		time.Sleep(100 * time.Millisecond)
	}
}
func (r *ring) prevDigitAnimated() {
	for i := 0; i < 3; i++ {
		(*r).prevRow()
		time.Sleep(100 * time.Millisecond)
	}
}

// builds a ring of rows
func newRing() *ring {
	d := newDigits()
	r := &ring{&sector{
		pos: 0,
		str: d[0],
	}}
	pred := r.sector

	for i := 1; i < 30; i++ {
		i := i
		// allocate new sector
		s := &sector{
			pos: i,
			str: d[i],
			// set predescessor as previous segment
			prev: pred,
		}
		// set predescessors next to point to new element
		(*pred).next = s

		if i == 29 {
			(*s).next = (*r).sector
			(*r.sector).prev = s
		}

		// let predescessor be the new element
		pred = s
	}
	// center ring on first element
	(*r).prevRow()

	// buffered de-/increase runs in the background, so that inc and dec
	// calls don't block (until the buffer is full) and
	// next/prevDigitAnimated calls run sequentially

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
func render(g *Game) {
	// painted, painted, painted… painted black
	termbox.Clear(0, 0)

	// COUNTER
	counter(g.state.points)

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
