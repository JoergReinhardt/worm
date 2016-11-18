package main

import (
	"fmt"
	"strings"
	"time"

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

type digit struct {
	abs  int
	ofl  bool
	next *digit
	buf  [3]string
}

func (d *digit) String() string {
	var str string
	if d.ofl {
		rbuf := strings.Split((*d.next).String(), "\n")
		str = fmt.Sprintf("%s %s\n%s %s\n%s %s\n",
			rbuf[0], (*d).buf[0], rbuf[1], (*d).buf[1], rbuf[2], (*d).buf[2])
	} else {
		str = fmt.Sprintf("%s\n%s\n%s\n", (*d).buf[0], (*d).buf[1], (*d).buf[2])
	}
	return str
}
func (d *digit) render(fn func(x, y int, c rune)) {
	w, _ := termbox.Size()
	// calculate size of boarders around message
	msg := strings.Split((*d).String(), "\n")
	for y, line := range msg {
		y, line := y, line
		lb := (w - len(line)) // left boarder
		for x, r := range line {
			x, r := x+lb, r
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
	(*d).buf[0] = digRows[d.abs*3]
	(*d).buf[1] = digRows[d.abs*3+1]
	(*d).buf[2] = digRows[d.abs*3+2]
}
func (d *digit) progDelayed() {
	delay := time.Millisecond * 500
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
		(*d).buf[0] = digRows[idx(d.abs, i, 0)]
		(*d).buf[1] = digRows[idx(d.abs, i, 1)]
		(*d).buf[2] = digRows[idx(d.abs, i, 2)]
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
	(*d).buf[0] = digRows[d.abs*3]
	(*d).buf[1] = digRows[d.abs*3+1]
	(*d).buf[2] = digRows[d.abs*3+2]
}
func (d *digit) regrDelayed() {
	delay := time.Millisecond * 500
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
		(*d).buf[0] = digRows[idx(d.abs, i, 0)]
		(*d).buf[1] = digRows[idx(d.abs, i, 1)]
		(*d).buf[2] = digRows[idx(d.abs, i, 2)]
		time.Sleep(delay)
	}
}

func newDigit() *digit {
	// initialize as digit zero
	return &digit{
		0,
		false,
		nil,
		[3]string{
			digRows[0],
			digRows[1],
			digRows[2],
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
	func() {
		for {
			select {
			case i := <-(*c).trig:
				if i > 0 {
					for i = i - 1; i > 0; {
						(*c.digit).progDelayed()
					}
				} else {
					for i = i + 1; i < 0; {
						(*c.digit).regrDelayed()
					}
				}
			}
		}
	}()
	return c
}
