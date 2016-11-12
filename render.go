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
	termbox.Clear(BLACK, BLACK)
	for y := 0; y <= g.state.width; y++ {
		for x := 0; x <= g.state.hight; x++ {
			switch {
			case (*g.cherry).picked(x, y): // set cherry
				termbox.SetCell(x, y, 'O', BLACK, RED)

			case (*g.worm).collides(x, y): // set worm segment
				termbox.SetCell(x, y, '#', BLACK, GREEN)

			default: // set all other cells to background
				termbox.SetCell(x, y, ' ', BLACK, BLACK)
			}
		}
	}
	termbox.Flush()
}
