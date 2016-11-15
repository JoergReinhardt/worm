package main

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
	"time"
)

func TestRing(t *testing.T) {
	r := newRing()
	var msg string

	for i := 0; i < 15; i++ {
		(*r).incDigit()
		msg = msg + r.stringDig()
	}
	t.Log(msg)
	msg = ""

	for i := 0; i < 57; i++ {
		(*r).advanceRow()
		msg = msg + (*r).stringRow() + "\n"
	}
	t.Log(msg)
	msg = ""

	for i := 0; i < 57; i++ {
		(*r).backupRow()
		msg = msg + (*r).stringRow() + "\n"
	}
	t.Log(msg)
	msg = ""

	for i := 0; i < 15; i++ {
		(*r).decDigit()
		msg = msg + r.stringDig()
	}
	t.Log(msg)
	msg = ""

	for i := 0; i < 15; i++ {
		(*r).incDigit()
		msg = msg + r.stringDig()
	}
	t.Log(msg)
	msg = ""

	(*r).nextDigitDelayed()
	go func() { (*r).nextDigitDelayed(); (*r).nextDigitDelayed() }()
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
func TestDigits(t *testing.T) {
	d := newDigit()
	t.Log(d.stringDig())
	for i := 0; i < 12; i++ {
		(*d).increase()
		t.Log(d.String())
	}
	t.Log(d.String())
	for i := 0; i < 100; i++ {
		(*d).increase()
		t.Log(d.String())
	}
	spew.Dump(d)
}
