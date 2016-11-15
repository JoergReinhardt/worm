package main

import (
	"testing"
)

func TestCounter(t *testing.T) {
	r := newRing()
	var msg string
	for i := 0; i < 100; i++ {
		(*r).sector = r.next
		msg = msg + r.sector.str + "\n"
	}
	t.Log(msg)
}
