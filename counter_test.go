package main

import (
	"testing"
	"time"
)

func TestDigit(t *testing.T) {
	d := newDigit()
	for i := 0; i < 10; i++ {
		(*d).progress()
		t.Log(d.String())
	}
	for i := 0; i < 10; i++ {
		(*d).progress()
		t.Log(d.String())
	}
	for i := 0; i < 10; i++ {
		(*d).regress()
		t.Log(d.String())
	}
	for i := 0; i < 10; i++ {
		(*d).regress()
		t.Log(d.String())
	}
	for i := 0; i < 10; i++ {
		(*d).progDelayed(time.Millisecond * 10)
		t.Log(d.String())
	}
	for i := 0; i < 10; i++ {
		(*d).regrDelayed(time.Millisecond * 10)
		t.Log(d.String())
	}
}
