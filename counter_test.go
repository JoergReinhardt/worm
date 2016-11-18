package main

import "testing"

func TestProgressDegress(t *testing.T) {
	d := newDigit()
	for i := 0; i < 20; i++ {
		(*d).progress()
		t.Log(d.String())
	}
	for i := 0; i < 20; i++ {
		(*d).regress()
		t.Log(d.String())
	}
}
func TestProgressDegressDelayed(t *testing.T) {
	d := newDigit()
	for i := 0; i < 20; i++ {
		(*d).progDelayed()
		t.Log(d.String())
	}
	for i := 0; i < 20; i++ {
		(*d).regrDelayed()
		t.Log(d.String())
	}
}
