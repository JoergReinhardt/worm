package main

import "testing"

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
}
