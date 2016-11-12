package main

import (
	//"bytes"
	//"fmt"
	"github.com/nsf/termbox-go"
)

const (
	WHITE termbox.Attribute = termbox.ColorWhite
	BLUE  termbox.Attribute = termbox.ColorBlue
	RED   termbox.Attribute = termbox.ColorRed
	GREEN termbox.Attribute = termbox.ColorGreen
	BLACK termbox.Attribute = termbox.ColorBlack
)

func render(g *Game) {
	// painted, painted, painted… painted black
	termbox.Clear(BLACK, BLACK)
	// render cherry
	termbox.SetCell((*g.cherry).x,
		(*g.cherry).y, 'O', BLACK, RED)

	// callback closes over SetCell with proper bg & fg, gets x &y by worms render method.
	var fn = func(x, y int) {
		termbox.SetCell(x, y, '#', BLACK, GREEN)
	}
	// renders through callback
	(*g.worm).render(fn)
	termbox.Flush()
}
