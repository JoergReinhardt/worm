package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
)

var digRows = []string{
	"┏━┓",
	"┃┃┃",
	"┗━┛",
	"╺┓ ",
	" ┃ ",
	"╺┻╸",
	"┏━┓",
	"┏━┛",
	"┗━╸",
	"┏━┓",
	"╺━┫",
	"┗━┛",
	"╻ ╻",
	"┗━┫",
	"  ╹",
	"┏━╸",
	"┗━┓",
	"┗━┛",
	"┏━┓",
	"┣━┓",
	"┗━┛",
	"┏━┓",
	"  ┃",
	"  ╹",
	"┏━┓",
	"┣━┫",
	"┗━┛",
	"┏━┓",
	"┗━┫",
	"┗━┛",
}

var delay = time.Millisecond * 500

type buf struct {
	*sync.RWMutex
	b [3]string
}
type digit struct {
	abs  int
	ofl  bool
	next *digit
	*buf
}

func (d *digit) String() string {
	var str string
	(*d.buf).Lock()
	defer (*d.buf).Unlock()
	if d.ofl {
		rbuf := strings.Split((*d.next).String(), "\n")
		str = fmt.Sprintf("%s %s\n%s %s\n%s %s\n",
			rbuf[0], (*d.buf).b[0], rbuf[1], (*d.buf).b[1], rbuf[2], (*d.buf).b[2])
	} else {
		str = fmt.Sprintf("%s\n%s\n%s\n", (*d.buf).b[0], (*d.buf).b[1], (*d.buf).b[2])
	}
	return str
}
func (d *digit) render(fn func(x, y int, c rune)) {
	w, _ := termbox.Size()
	// calculate size of boarders around message
	msg := strings.Split((*d).String(), "\n")
	lb := (w - len(msg[0])) // left boarder
	for y, line := range msg {
		y, line := y, line
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
func (d *digit) progress() {
	if d.abs == 9 {
		(*d).abs = 0
		if !d.ofl {
			(*d).next = newDigit()
			(*d).ofl = true
		}
		(*d).next.progress()
	} else {
		(*d).abs = d.abs + 1
	}
	(*d.buf).b[0] = digRows[d.abs*3]
	(*d.buf).b[1] = digRows[d.abs*3+1]
	(*d.buf).b[2] = digRows[d.abs*3+2]
}
func (d *digit) progDelayed() {
	if d.abs == 9 {
		(*d).abs = 0
		if !d.ofl {
			(*d).next = newDigit()
			(*d).ofl = true
		}
		(*d).next.progDelayed()
	} else {
		(*d).abs = d.abs + 1
	}
	var idx = func(abs, idx, rcn int) int {
		// wrap over and underflow
		i := 3 * abs
		if i != 0 {
			i = i % 30
		}
		return i + rcn
	}
	for i := 0; i < 3; i++ {
		i := i
		(*d.buf).Lock()
		(*d.buf).b[0] = digRows[idx(d.abs, i, 0)]
		(*d.buf).b[1] = digRows[idx(d.abs, i, 1)]
		(*d.buf).b[2] = digRows[idx(d.abs, i, 2)]
		(*d.buf).Unlock()
		time.Sleep(delay)
	}
}
func (d *digit) regress() {
	if d.abs == 0 {
		(*d).abs = 9
		if d.ofl {
			(*d).next.regress()
		}
	} else {
		(*d).abs = d.abs - 1
	}
	(*d.buf).b[0] = digRows[d.abs*3]
	(*d.buf).b[1] = digRows[d.abs*3+1]
	(*d.buf).b[2] = digRows[d.abs*3+2]
}
func (d *digit) regrDelayed() {
	if d.abs == 0 {
		(*d).abs = 9
		if d.ofl {
			(*d).next.regress()
		}
	} else {
		(*d).abs = d.abs - 1
	}
	var idx = func(abs, idx, rcn int) int {
		// wrap over and underflow
		i := 3 * abs
		if i != 0 {
			i = i % 30
		}
		return i + rcn
	}
	for i := 3; i > 0; i-- {
		i := i
		(*d.buf).Lock()
		(*d.buf).b[0] = digRows[idx(d.abs, i, 0)]
		(*d.buf).b[1] = digRows[idx(d.abs, i, 1)]
		(*d.buf).b[2] = digRows[idx(d.abs, i, 2)]
		(*d.buf).Unlock()
		time.Sleep(delay)
	}
}

func newDigit() *digit {
	// initialize as digit zero
	return &digit{
		0,
		false,
		nil,
		&buf{
			&sync.RWMutex{},
			[3]string{
				digRows[0],
				digRows[1],
				digRows[2],
			},
		},
	}
}

type counter struct {
	*digit
	trig chan int
}

func (c *counter) inc(i int) { (*c).trig <- i }
func newCounter() *counter {
	trig := make(chan int, 16)
	c := &counter{newDigit(), trig}
	go func() {
		for {
			select {
			case i := <-(*c).trig:
				if i > 1 {
					for j := 0; j < i; j++ {
						fmt.Printf("trig %d ", i)
						(*c).trig <- 1
					}
				}
				if i < -1 {
					for j := 0; j < i; j-- {
						(*c).trig <- -1
					}
				}
				if i == 1 {
					(*c.digit).progDelayed()
					i = i - 1
				}
				if i == -1 {
					(*c.digit).regrDelayed()
					i = i + 1
				}
			}
		}
	}()
	return c
}
