package main

import (
	"math/rand"
	"time"
)

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
	}
	return false
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
	speed      time.Duration
	eventState gameStat
	direction  dir
	size       func() (x, y int)
	rand       func() (x, y int)
}

// newState()
// to allocate a new state struct, two functional arguments need to be passed,
// to retrieve current Boardsize and generate random positions within board
// dimensions, in order to wrap the board according to it's size and relocate
// the cherry, if it got picked.
func newState(sizeFn func() (x, y int), randFn func() (x, y int)) *state {
	return &state{
		speed:      250 * time.Millisecond,
		eventState: INIT,
		direction:  UP,
		size:       sizeFn,
		rand:       randFn,
	}
}

// Game the game struct holds all game elements and its current state.
type game struct {
	*state
	*cherry
	*worm
}

// allocate a new game
func newGame(sizeFn func() (x, y int)) *game {
	// initialize random number generation
	rand.Seed(time.Now().Unix())
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
	return &game{s, &cherry{x, y}, newWorm(s)}
}

// reset game, by replacing it with a new game
func (g *game) reset() {
	*g = *newGame(g.state.size)
	points = newDigit()
}

// bend board in third dimension to become a continuous manifold ;)
func (g *game) wrapBoard(xi, yi int) (xo, yo int) {
	xo, yo = xi, yi
	w, h := g.state.size()
	if xi < 0 { // wrap left boarder to right
		xo = w - 1
	}
	if xi >= w { // wrap right boarder to left
		xo = 0
	}
	if yi < 0 { // wrap upper boarder to lower
		yo = h - 1
	}
	if yi >= h { //wrap lower boarder to upper
		yo = 0
	}
	return xo, yo
}

// does one worm move and all neccessary changes that follow by the new state.
func (g *game) play() {

	// wrap board to continuoum, get worms final x & y regarding the
	// current board size
	x, y := (*g).wrapBoard((*g.worm).predNextPos(g.state.direction))

	// IF SELF-COLLDISION → GAME OVER
	if (*g.worm).collides(x, y) {
		(*g.state).eventState = GAME_OVER
		return
	}

	// cherry possibly needs to be relocatet. In case boardsize changed
	// cherry coordinates may turn out to not be on the board anymore
	(*g.cherry).x, (*g.cherry).y = (*g).wrapBoard((*g.cherry).x, (*g.cherry).y)

	// CHECK IF CHERRY GOT PICKED
	if (*g.cherry).picked(x, y) {
		// grow worm
		(*g.worm).grow()
		// relocate cherry
		(*g.cherry).pop(g.state.rand())
		// increase worm speed by 10%
		(*g.state).speed = (g.state.speed / 10) * 9
		// raise points by one
		(*points).increase()
	}
	// MOVE ON TO NEXT POSITION
	(*g.worm).move(x, y, g.state.direction)
}
