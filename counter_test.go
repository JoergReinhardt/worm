package main

import (
	"testing"
	"time"
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

	for i := 0; i < 15; i++ {
		(*r).nextDigit()
		msg = msg + r.stringDig()
	}
	t.Log(msg)
	msg = ""

	(*r).nextDigitAnimated()
	go func() { (*r).nextDigitAnimated(); (*r).nextDigitAnimated() }()
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	t.Log(r.stringDig())
	time.Sleep(50 * time.Millisecond)
	msg = ""
}
