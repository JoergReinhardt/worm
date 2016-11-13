package main

import (
	"fmt"
	"time"
)

type dir uint8

const (
	UP   dir = 0
	DOWN dir = 1 << iota
	LEFT
	RIGHT
)

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

// takes a closure over termbox SetCell and passes the coordinates of each
// segment.
func (s *segment) render(fn func(x, y int)) {
	fn(s.x, s.y)
	if s.tail {
		return
	} else {
		(*s.next).render(fn)
	}
}

type worm struct {
	*segment
}

func (w worm) predict(s *state) (x, y int) {
	ox, oy := w.segment.x, w.segment.y

	// set new position for this segment
	switch s.move {
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
func newWorm(s *state) *worm {
	x, y := s.size()
	return &worm{
		segment: &segment{
			x:    x / 2,
			y:    y / 2,
			tail: true,
			next: nil,
		},
	}
}

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

type state struct {
	speed time.Duration
	stat  GameStat
	move  dir
	size  func() (x, y int)
}

//go:generate stringer -type GameStat
type GameStat uint8

const (
	INIT GameStat = 0
	RUN  GameStat = 1 << iota
	PAUSE
	GAME_OVER
)

type Game struct {
	*state
	*cherry
	*worm
}

func NewGame(sfn func() (x, y int)) *Game {
	s := &state{250 * time.Millisecond, 0, UP, sfn}
	return &Game{s, &cherry{20, 20}, newWorm(s)}
}

func (g *Game) wrapBoard(xi, yi int) (xo, yo int) {
	xo, yo = xi, yi
	w, h := g.state.size()
	if xi < 0 {
		xo = w - 1
	}
	if xi == w {
		xo = 0
	}
	if yi < 0 {
		yo = h - 1
	}
	if yi == h {
		yo = 0
	}
	return xo, yo
}
func (g *Game) move() {
	// predict next positiom
	x, y := (*g.worm).predict(g.state)
	x, y = (*g).wrapBoard(x, y)
	// if next position on cherry
	if (*g.cherry).picked(x, y) {
		// grow worm
		(*g.worm).grow()
	}
	// GAME OVER
	if (*g.worm).collides(x, y) {
		(*g.state).stat = GAME_OVER
		return
	}
	(*g.worm).move(x, y)
}
