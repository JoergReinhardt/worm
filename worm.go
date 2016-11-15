package main

// the char a segment is showing, depends on the position of previous and next
// segments position, which opens the opportunety for a map, utilising an
// anonymous two field struct as it's key. (just because go can do that)
var segChars = map[struct{ p, n dir }]rune{
	struct{ p, n dir }{UP, DOWN}:    '┃',
	struct{ p, n dir }{DOWN, UP}:    '┃',
	struct{ p, n dir }{LEFT, RIGHT}: '━',
	struct{ p, n dir }{RIGHT, LEFT}: '━',
	struct{ p, n dir }{LEFT, UP}:    '┛',
	struct{ p, n dir }{UP, LEFT}:    '┛',
	struct{ p, n dir }{LEFT, DOWN}:  '┓',
	struct{ p, n dir }{DOWN, LEFT}:  '┓',
	struct{ p, n dir }{RIGHT, UP}:   '┗',
	struct{ p, n dir }{UP, RIGHT}:   '┗',
	struct{ p, n dir }{RIGHT, DOWN}: '┏',
	struct{ p, n dir }{DOWN, RIGHT}: '┏',
}

// each of worms segments holds it's own posisiton, the char to render, a
// boolflag if there is a next element and a pointer to that element (nil
// pointer if not present, hence the bool flag for convienience.)
type segment struct {
	x, y int
	char rune
	tail bool
	next *segment
}

// checks if passed coordinates collide with the coordinates of this segment
// and all it"s children.
func (s *segment) collides(x, y int) bool {
	// when coordinates are identical, collide:
	if s.x == x && s.y == y {
		return true
	}
	if s.tail { // otherwise, return no collision
		return false
	} else { // when not tail, check for tail collision
		return s.next.collides(x, y)
	}
}

// allocates a new segment and assignes it to the 'next' field, if this is the
// tail, otherwise delegates that task to next element.
func (s *segment) grow() {
	if s.tail { // try to add new tail element to this segment
		// if this segment was tail before, it is not any longer.
		(*s).tail = false
		// initialize new tail segment at current position
		(*s).next = &segment{s.x, s.y, s.char, true, nil}
	} else { // if not tail, deligate to next segment
		(*s.next).grow()
	}
}

// takes a closure over termbox SetCell with preset bf & fg color as argument
// and calls it, passing elements current and all it's childs coordinates and
// chars recursesively.
func (s *segment) render(fn func(x, y int, c rune)) {
	fn(s.x, s.y, s.char)
	if !s.tail {
		(*s.next).render(fn)
	}
}

// returns the direction this element is located at, relative to the given
// position
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

// move to passed position, choose char, based on previous and next elements
// relative position and drag all childs along recursively.
func (s *segment) move(x, y, px, py int) {
	// safe elements current position
	nx, ny := (*s).x, (*s).y
	// move to new position, passed by caller
	(*s).x, (*s).y = x, y
	if s.tail { // if this is the tail element:
		// get tail char, based on previous elements relative position
		switch s.RelPos(px, py) {
		case UP:
			(*s).char = '┆'
		case DOWN:
			(*s).char = '┆'
		case LEFT:
			(*s).char = '┄'
		case RIGHT:
			(*s).char = '┄'
		}
	} else { // if this is an ordinary segment:
		// allocate anonymous struct to hold two relative positions
		var relPos = struct {
			p dir
			n dir
		}{
			s.RelPos(px, py),
			s.RelPos(nx, ny),
		}
		// get new char based on relative previous and next positions
		(*s).char = segChars[relPos]
		// pass old position as new position for the next element. pass
		// self to determin relative position.
		(*s.next).move(nx, ny, x, y)
	}
}

////////////////////////////////////////////////////////////////////////////////
// worm is a singly linked list of segments. It's 'head' can predict it's next
// position, according to current direction of movemen.
type worm struct {
	*segment
}

func (w *worm) move(x, y int, d dir) {
	// safe current position as next elements future position
	nx, ny := w.segment.x, w.segment.y
	// move element to new position (passed by caller)
	(*w.segment).x, (*w.segment).y = x, y
	// determine head char based on current direction of movement
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

	if w.tail { // if head also happens to be tail, render headchar and
		// be done with it
		return
	} else { // otherwise pass old position and pointer to self on to the
		// next segments move method
		(*w.segment.next).move(nx, ny, x, y)
	}
}
func (w worm) predNextPos(d dir) (x, y int) {
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
			char: ' ',
			tail: true,
			next: nil,
		},
	}
}
