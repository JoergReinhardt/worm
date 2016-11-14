package main

import (
	"fmt"
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

// allocates a new segment and allocates it to 'next', if its the tail,
// otherwise delegates that task to next element recursively.
func (s *segment) grow() {
	if s.tail { // try to add new tail element to this segment
		(*s).tail = false // if this segment is tail, it is not any
		// longer initialize new tail segment at current position
		(*s).next = &segment{s.x, s.y, s.char, true, nil}
	} else { // if not tail, deligate to next segment
		(*s.next).grow()
	}
}

// takes a closure over termbox SetCell with preset bf & fg color and char as
// argument and calls it, passing elements current and all it's childs
// coordinates recursesively.
func (s *segment) render(fn func(x, y int, c rune)) {
	fn(s.x, s.y, s.char)
	if s.tail {
		return
	} else {
		(*s.next).render(fn)
	}
}

// returns the direction this element is located relative to the given
func (s segment) RelPos(x, y int) (d dir) {
	switch {
	case s.y < y:
		d = DOWN // ABOVE
	case s.y > y:
		d = UP // BEYOND
	case s.x < x:
		d = RIGHT // LEFT-OF
	case s.x > x:
		d = LEFT // RIGHT-OF
	}
	return d
}

var segChars = map[struct{ p, n dir }]rune{
	struct{ p, n dir }{UP, DOWN}:    '┃',
	struct{ p, n dir }{DOWN, UP}:    '┃',
	struct{ p, n dir }{LEFT, RIGHT}: '━',
	struct{ p, n dir }{RIGHT, LEFT}: '━',
	struct{ p, n dir }{LEFT, UP}:    '┛',
	struct{ p, n dir }{LEFT, DOWN}:  '┓',
	struct{ p, n dir }{RIGHT, UP}:   '┗',
	struct{ p, n dir }{RIGHT, DOWN}: '┏',
	struct{ p, n dir }{UP, LEFT}:    '┛',
	struct{ p, n dir }{UP, RIGHT}:   '┗',
	struct{ p, n dir }{DOWN, LEFT}:  '┓',
	struct{ p, n dir }{DOWN, RIGHT}: '┏',
}

// move to passed position and drag all childs along recursively.
func (s *segment) move(x, y int, prev *segment) {
	// safe elements current position
	px, py := prev.x, prev.y
	nx, ny := (*s).x, (*s).y
	// move tail to the passed position
	(*s).x, (*s).y = x, y
	if s.tail {
		// get char, based on previous elements direction
		switch prev.RelPos(x, y) {
		case UP:
			(*s).char = ','
		case DOWN:
			(*s).char = '\''
		case LEFT:
			(*s).char = '-'
		case RIGHT:
			(*s).char = '~'
		}
	} else { // recursively call move for tail elements
		// allocate anonymous struct to hold two relative positions
		var relPos struct{ p, n dir }
		// determin position relative to previous elements position
		// (passed by caller)
		relPos.p = s.RelPos(px, py)
		// determin position relative to next elements position (identical to
		// old position)
		relPos.n = s.RelPos(nx, ny)
		// get new char based on relative previous and next positions
		(*s).char = segChars[relPos]
		// pass old position as new position for the next element. pass
		// self to determin relative position.
		(*s.next).move(nx, ny, s)
	}
}

// worm is a singly linked list of segments. It's 'head' can predict it's next
// position, according to current direction of movemen.
type worm struct {
	queue [2]*segment
	*segment
}

func (w *worm) move(x, y int, d dir) {
	// safe current position as next elements future position
	nx, ny := w.segment.x, w.segment.y
	// move element to new position
	(*w.segment).x, (*w.segment).y = x, y
	// determine char based on current direction
	switch d {
	case UP:
		(*w.segment).char = '^'
	case DOWN:
		(*w.segment).char = 'v'
	case LEFT:
		(*w.segment).char = '<'
	case RIGHT:
		(*w.segment).char = '>'
	}

	if w.tail {
		return
	} else {
		(*w.segment.next).move(nx, ny, w.segment)
	}
}
func (w worm) predict(d dir) (x, y int) {
	// get current (old) position
	ox, oy := w.segment.x, w.segment.y
	// predict next position regarding current direction and position
	switch d {
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
