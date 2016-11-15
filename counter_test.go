package main

import (
	"testing"
)

func TestCounter(t *testing.T) {
	r := newRing()
	var msg string
	for i := 0; i < 33; i++ {
		(*r).nextRow()
		msg = msg + (*r).stringRow() + "\n"
	}
	t.Log(msg)
	for i := 0; i < 33; i++ {
		(*r).prevRow()
		msg = msg + (*r).stringRow() + "\n"
	}
	t.Log(msg)
}
