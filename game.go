package main

import (
	"fmt"
)

type dir uint8

const (
	UP   dir = 0
	DOWN dir = 1 << iota
	LEFT
	RIGHT
)

var DIRECTION = UP
var WIDTH = 79
var HIGHT = 79

type segment struct {
	x, y int
	tail bool
	next *segment
}

func (s segment) String() string {
	var str = fmt.Sprint(s.x) + ", " + fmt.Sprint(s.y) + "\n"
	if !s.tail {
		str = str + s.next.String()
	}
	return str
}
func (s segment) len() int {
	if s.tail {
		return 1
	} else { // calculate length recursively
		return s.next.len() + 1
	}
}
func (s *segment) grow() {
	if s.tail { // add new tail element
		(*s).tail = false // this segment is not tail any longer
		// initialize next segment at current position
		(*s).next = &segment{s.x, s.y, true, nil}
	} else { // deligate to next node
		(*s.next).grow()
	}
}
func (s *segment) move(x, y int) {
	// safe old position
	ox, oy := (*s).x, (*s).y
	// move to new position
	(*s).x, (*s).y = x, y
	if !s.tail { // recursively call move for tail elements
		(*s.next).move(ox, oy)
	}
}
func (s *segment) collides(x, y int) bool {
	// when coordinates are identical, collide:
	if s.x == x && s.y == y {
		return true
	}
	if !s.tail { // when not tail, check for tail collision
		return s.next.collides(x, y)
	} else { // otherwise, return no collision
		return false
	}
}

type worm struct {
	*segment
}

func (w worm) predict() (x, y int) {
	ox, oy := w.segment.x, w.segment.y

	// set new position for this segment
	switch DIRECTION {
	case UP:
		y = oy - 1
		x = ox
	case DOWN:
		y = oy + 1
		x = ox
	case LEFT:
		x = ox - 1
		y = oy
	case RIGHT:
		x = ox + 1
		y = oy
	}
	return x, y
}

// set to middle of field
func newWorm() *worm {
	return &worm{
		segment: &segment{
			x:    WIDTH / 2,
			y:    HIGHT / 2,
			tail: true,
			next: nil,
		},
	}
}

///////////////////////////////////////////////////////
type color uint

const (
	BLANK color = 0
	WORM  color = 1 << iota
	CHERRY
)

///////////////////////////////////////////////////////
type cherry struct {
	x int
	y int
}

func (c cherry) picked(x, y int) bool {
	if c.x == x && c.y == y {
		return true
	} else {
		return false
	}
}

//go:generate stringer -type State
type State uint

const (
	INIT State = 0
	RUN  State = 1 << iota
	PAUSE
	GAME_OVER
)

type Game struct {
	State
	*cherry
	*worm
}

func NewGame() *Game {
	return &Game{0, &cherry{20, 20}, newWorm()}
}

func (g *Game) wrap(xi, yi int) (xo, yo int) {
	xo, yo = xi, yi
	if xi < 0 {
		xo = WIDTH - 1
	}
	if xi == WIDTH {
		xo = 0
	}
	if yi < 0 {
		yo = HIGHT - 1
	}
	if yi == HIGHT {
		yo = 0
	}
	return xo, yo
}
func (g *Game) move() {
	// predict next positiom
	x, y := (*g.worm).predict()
	x, y = (*g).wrap(x, y)
	// if next position on cherry
	if (*g.cherry).picked(x, y) {
		// grow worm
		(*g.worm).grow()
	}
	// GAME OVER
	if (*g.worm).collides(x, y) {
		(*g).State = GAME_OVER
		return
	}
	(*g.worm).move(x, y)
}
