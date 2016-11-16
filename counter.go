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
	buf [3]string
}

func (d *digit) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n", d.buf[0], d.buf[1], d.buf[2])
}
func (d *digit) progress() {
	if d.abs == 9 {
		(*d).abs = 0
	} else {
		(*d).abs = d.abs + 1
	}
	(*d).buf[0] = digRows[d.abs*3]
	(*d).buf[1] = digRows[d.abs*3+1]
	(*d).buf[2] = digRows[d.abs*3+2]
}
func (d *digit) progDelayed(delay time.Duration) {
	if d.abs == 9 {
		(*d).abs = 0
	} else {
		(*d).abs = d.abs + 1
	}
	for i := 0; i < 3; i++ {
		i := i
		(*d).buf[0] = digRows[i+d.abs*3]
		(*d).buf[1] = digRows[i+d.abs*3+1]
		(*d).buf[2] = digRows[i+d.abs*3+2]
		time.Sleep(delay)
	}
	(*d).abs = d.abs + 1
}
func (d *digit) regress() {
	if d.abs == 0 {
		(*d).abs = 9
	} else {
		(*d).abs = d.abs - 1
	}
	(*d).buf[0] = digRows[d.abs*3]
	(*d).buf[1] = digRows[d.abs*3+1]
	(*d).buf[2] = digRows[d.abs*3+2]
}
func (d *digit) regDelayed(delay time.Duration) {
	for i := 2; i > 0; i-- {
		i := i
		(*d).buf[0] = digRows[i+d.abs%30]
		(*d).buf[1] = digRows[i+d.abs%30+1]
		(*d).buf[2] = digRows[i+d.abs%30+2]
		time.Sleep(delay)
	}
	(*d).abs = d.abs - 1
}

func newDigit() *digit {
	// initialize as digit zero
	return &digit{
		0,
		[3]string{
			digRows[0],
			digRows[1],
			digRows[2],
		},
	}
}
