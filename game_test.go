package main

import (
	"testing"
)

func TestSegment(t *testing.T) {
	g := NewGame()
	t.Log(g.worm)
	DIRECTION = UP
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	DIRECTION = LEFT
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	t.Log(g)
	(*g).cherry.x = 25
	(*g).cherry.y = 25
	DIRECTION = DOWN
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	t.Log(g.worm)
	DIRECTION = RIGHT
	g.move()
	g.move()
	g.move()
	g.move()
	g.move()
	t.Log(g)
	DIRECTION = LEFT
	g.move()
	t.Log(g.State)
	t.Log(g.worm)
}
