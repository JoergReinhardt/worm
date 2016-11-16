package main

import (
	"fmt"
	"time"
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
	abs int
	ofl bool
	buf [3]string
}

func (d *digit) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n", d.buf[0], d.buf[1], d.buf[2])
}
func (d *digit) progress() {
	if d.abs == 9 {
		(*d).abs = 0
		(*d).ofl = true
	} else {
		(*d).abs = d.abs + 1
	}
	(*d).buf[0] = digRows[d.abs*3]
	(*d).buf[1] = digRows[d.abs*3+1]
	(*d).buf[2] = digRows[d.abs*3+2]
}
func (d *digit) progDelayed(delay time.Duration) {
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
	if d.abs > 9 {
		(*d).abs = 0
		(*d).ofl = true
	} else {
		(*d).abs = d.abs + 1
	}
}
func (d *digit) regress() {
	if d.abs == 0 {
		(*d).abs = 9
		(*d).ofl = true
	} else {
		(*d).abs = d.abs - 1
	}
	(*d).buf[0] = digRows[d.abs*3]
	(*d).buf[1] = digRows[d.abs*3+1]
	(*d).buf[2] = digRows[d.abs*3+2]
}
func (d *digit) regrDelayed(delay time.Duration) {
	var idx = func(abs, idx, rcn int) int {
		// wrap over and underflow
		i := 3 * abs
		if i != 0 {
			i = i % 30
		}
		return i + rcn
	}
	if d.abs == 0 {
		(*d).abs = 9
		(*d).ofl = true
	} else {
		(*d).abs = d.abs - 1
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
		[3]string{
			digRows[0],
			digRows[1],
			digRows[2],
		},
	}
}
