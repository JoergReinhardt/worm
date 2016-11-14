package main

import (
	"fmt"
	"math/rand"
	"time"
)

type dir uint8

const (
	UP   dir = 0
	DOWN dir = 1 << iota
	LEFT
	RIGHT
)

// each of worms segments holds it's own posisiton, a boolflag if there is a
// next element and a pointer to that element (nil pointer if not present,
// hence the bool flag for convienience.)
type segment struct {
	x, y int
	char rune
	tail bool
	next *segment
}

// string and len as helpers during debug phase get not used anymore
func (s segment) String() string {
	var str = fmt.Sprint(s.x) + ", " + fmt.Sprint(s.y) + "\n"
	if !s.tail {
		str = str + s.next.String()
	}
	return str
}

// string and len as helpers during debug phase get not used anymore
func (s segment) len() int {
	if s.tail {
		return 1
	} else { // calculate length recursively
		return s.next.len() + 1
	}
}

// allocates a new segment and allocates it to 'next', if its the tail,
// otherwise delegates that task to next element recursively.
func (s *segment) grow() {
	if s.tail { // try to add new tail element to this segment
		(*s).tail = false // if this segment is tail, it is not any
		// longer initialize new tail segment at current position
		(*s).next = &segment{s.x, s.y, '~', true, nil}
	} else { // if not tail, deligate to next segment
		(*s.next).grow()
	}
}

// move to passed position and drag all childs along recursively.
func (s *segment) move(x, y int, char rune) {
	// safe old position
	ox, oy := s.x, s.y
	oc := s.char
	// move to new position
	(*s).x, (*s).y = x, y
	if !s.tail { // recursively call move for tail elements
		switch oc {
		case '^', 'v', 'A', 'Y':
			(*s.next).move(ox, oy, ':')
		case '<', '>', '(', '[', ')', ']':
			(*s.next).move(ox, oy, '~')
		}
	}
	switch oc {
	case '^', 'A':
		wr(w.x, w.y, 'O')
	case 'v', 'Y':
		wr(w.x, w.y, 'o')
	case 'O':
		wr(w.x, w.y, 'A')
	case 'o':
		wr(w.x, w.y, 'Y')
	case '>', ')':
		wr(w.x, w.y, ']')
	case '<', '(':
		wr(w.x, w.y, '[')
	case ']':
		wr(w.x, w.y, '(')
	case '[':
		wr(w.x, w.y, '(')
	}
}

// checks if passed coordinates collide with the coordinates of this segment
// and all it"s childs.
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

// takes a closure over termbox SetCell with preset bf & fg color and char as
// argument and calls it, passing elements current and all it's childs
// coordinates recursesively.
func (s *segment) render(wr func(x, y int, char rune), oc rune) {
}

// worm is a singly linked list of segments. It's 'head' can predict it's next
// position, according to current direction of movemen.
type worm struct {
	*segment
}

func (w *worm) render(wr func(x, y int, char rune), d dir) {
	var char rune
	switch d {
	case UP:
		wr(w.x, w.y, '^')
	case DOWN:
		wr(w.x, w.y, 'v')
	case LEFT:
		wr(w.x, w.y, '<')
	case RIGHT:
		wr(w.x, w.y, '>')
	}

}

func (w worm) predict(s *state) (x, y int) {
	// get current (old) position
	ox, oy := w.segment.x, w.segment.y
	// predict next position regarding current direction and position
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
	// return predicted next position
	return x, y
}

// center new worm at board
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
// nice jucy relocateable cherry as food for our worm
type cherry struct {
	x int
	y int
}

// indicates, if coordinates identical to cherrys coordinates, hence 'if the
// cherry is picked'
func (c cherry) picked(x, y int) bool {
	if c.x == x && c.y == y {
		return true
	} else {
		return false
	}
}

// pops cherry at a new location
func (c *cherry) pop(x, y int) {
	(*c).x, (*c).y = x, y
}

///////////////////////////////////////////////////////
/// STATE
// gets passed araound to cooperate with controller and render functions on
// progressing the game state.
type state struct {
	speed time.Duration
	stat  GameStat
	move  dir
	size  func() (x, y int)
	rand  func() (x, y int)
}

// newState()
// to allocate a new state struct, two functional arguments need to be passed,
// to retrieve current Boardsize and generate random positions within board
// dimensions, in order to wrap the board according to it's size and relocate
// the cherry, if it got picked.
func newState(sizeFn func() (x, y int), randFn func() (x, y int)) *state {
	return &state{
		250 * time.Millisecond,
		INIT,
		UP,
		sizeFn,
		randFn,
	}
}

// current Game state flag
//go:generate stringer -type GameStat
type GameStat uint8

const (
	INIT GameStat = 0
	RUN  GameStat = 1 << iota
	PAUSE
	GAME_OVER
)

// the game struct holds all game elements and its current state.
type Game struct {
	*state
	*cherry
	*worm
}

// allocate a new game
func NewGame(sizeFn func() (x, y int)) *Game {
	// initialize random number generation
	rand.Seed(23)
	// build a closure to generate random x & y, within the range of width
	// & hight of the board, given by the size functions Output, which this
	// closure closes over.
	randFn := func() (x, y int) {
		w, h := sizeFn()
		return rand.Intn(w), rand.Intn(h)
	}
	x, y := randFn()              // initial cherry position
	s := newState(sizeFn, randFn) // allocate new state
	// return game with all its components
	return &Game{s, &cherry{x, y}, newWorm(s)}
}

// bend board in third dimension to become a continuous manifold ;)
func (g *Game) wrapBoard(xi, yi int) (xo, yo int) {
	xo, yo = xi, yi
	w, h := g.state.size()
	if xi < 0 { // wrap left boarder to right
		xo = w - 1
	}
	if xi == w { // wrap right boarder to left
		xo = 0
	}
	if yi < 0 { // wrap upper boarder to lower
		yo = h - 1
	}
	if yi == h { //wrap lower boarder to upper
		yo = 0
	}
	return xo, yo
}

// does one worm move and all neccessary changes that follow by the new state.
func (g *Game) play() {
	// predict worms next positiom
	x, y := (*g.worm).predict(g.state)
	// wrap board to continuoum
	x, y = (*g).wrapBoard(x, y)
	// IF CHERRY GOT PICKED
	if (*g.cherry).picked(x, y) {
		// grow worm
		(*g.worm).grow()
		// relocate cherry
		(*g.cherry).pop(g.state.rand())
		// raise worm speed by 10%
		(*g.state).speed = (g.state.speed / 10) * 9
	}
	// IF SELF-COLLDISION → GAME OVER
	if (*g.worm).collides(x, y) {
		(*g.state).stat = GAME_OVER
		return
	}
	// MOVE TO NEXT POSITION
	(*g.worm).move(x, y)
}
