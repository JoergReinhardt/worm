package main

import (
	"testing"
)

func TestCounter(t *testing.T) {
	r := newRing()
	var msg string

	for i := 0; i < 15; i++ {
		(*r).nextDigit()
		msg = msg + r.stringDig()
	}
	t.Log(msg)
	msg = ""

	for i := 0; i < 57; i++ {
		(*r).nextRow()
		msg = msg + (*r).stringRow() + "\n"
	}
	t.Log(msg)
	msg = ""

	for i := 0; i < 57; i++ {
		(*r).prevRow()
		msg = msg + (*r).stringRow() + "\n"
	}
	t.Log(msg)
	msg = ""

	for i := 0; i < 15; i++ {
		(*r).prevDigit()
		msg = msg + r.stringDig()
	}
	t.Log(msg)
	msg = ""
}
