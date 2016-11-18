package main

import (
	"fmt"
	"testing"
)

type termBoxMock struct {
	buf [][]rune
}

func newTermBoxMock() *termBoxMock {
	return &termBoxMock{[][]rune{[]rune{}}}
}
func (tb *termBoxMock) SetCell(x, y int, r rune) {

	fmt.Print(r)
	// bring buffer to appropriate size
	if len((*tb).buf) <= y {
		(*tb).buf = append((*tb).buf, []rune{})
	}
	// set cell
	fmt.Print(r)
	(*tb).buf[x] = append(tb.buf[x], r)
}
func (tb termBoxMock) String() string {
	var str string
	for _, line := range tb.buf {
		line := line
		fmt.Sprintf("%s\n", line)
	}
	return str
}
func TestInitScreen(t *testing.T) {
	tbm := newTermBoxMock()
	initScreen.render((*tbm).SetCell)
	t.Log(tbm.buf)
}
func TestCounter(t *testing.T) {
	tbm := newTermBoxMock()
	cnt := newCounter()
	for i := 0; i < 23; i++ {
		(*cnt).inc(1)
	}
	cnt.render((*tbm).SetCell)

	t.Log(tbm.buf)
}
