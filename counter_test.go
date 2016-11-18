package main

import (
	"testing"
	"time"
)

func TestCounterString(t *testing.T) {
	cnt := newCounter()
	(*cnt).inc(5)

	time.Sleep(5 * time.Second)

	t.Log(cnt.String())
}
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
