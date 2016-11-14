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
	termbox.SetCell((*g.cherry).x, (*g.cherry).y, 'O', RED, BLACK)

	// wraps termbox.SetCell and presetes colors.
	var renderWorm = func(x, y int, char rune) {
		termbox.SetCell(x, y, char, WHITE, BLACK)
	}

	// pass the worm renderer tto worms render method
	(*g.worm).render(renderWorm, g.state.move)

	// flush render buffer
	termbox.Flush()

}
