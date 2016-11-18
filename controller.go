package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

const (
	RED   termbox.Attribute = termbox.ColorRed
	GREEN termbox.Attribute = termbox.ColorGreen
	BLACK termbox.Attribute = termbox.ColorBlack
	WHITE termbox.Attribute = termbox.ColorWhite
)

type dir uint8

const (
	UP   dir = 0
	DOWN dir = 1 << iota
	LEFT
	RIGHT
)

type gameStat uint8

const (
	INIT gameStat = 0
	RUN  gameStat = 1 << iota
	PAUSE
	GAME_OVER
)

func run() { // runs the animation and input event cycles

	//////////////////////////////////////////////////////////////////////////////////////
	/// GAME SETUP
	// Allocate new game, pass termbox size as functional argument
	g := newGame(termbox.Size)
	// game loop progresses the game in the background
	go gameController(g)

	//////////////////////////////////////////////////////////////////////////////////////
	/// EVENT QUEUE SETUP
	// the event queue yields keyboard events
	eventQueue := make(chan termbox.Event)
	// dplete event queue in the background
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	//////////////////////////////////////////////////////////////////////////////////////
	/// EVENT LOOP
	// event retrieval and render loop. Blocks until game has stopped.
	for {
		ev := <-eventQueue
		if ev.Type == termbox.EventKey {
			switch {
			// set direction for the next move
			case ev.Ch == 's':
				// if on init screen, run, when s is pressed,
				// else ignore
				if g.state.eventState == INIT {
					(*g.state).eventState = RUN
				}
			case ev.Ch == 'p':
				// toggle game state between pause & run
				if g.state.eventState == PAUSE {
					(*g.state).eventState = RUN
				} else {
					(*g.state).eventState = PAUSE
				}
			case ev.Key == termbox.KeyArrowUp || ev.Ch == 'k':
				g.state.direction = UP
			case ev.Key == termbox.KeyArrowDown || ev.Ch == 'j':
				g.state.direction = DOWN
			case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'h':
				g.state.direction = LEFT
			case ev.Key == termbox.KeyArrowRight || ev.Ch == 'l':
				g.state.direction = RIGHT
			case ev.Key == termbox.KeyEsc || ev.Ch == 'q':
				return
			}
		}

		if g.state.eventState == GAME_OVER {
			(*g).reset()
			(*g.state).eventState = INIT

		}
		(*g).render()

		// sleep til next animation cycle
		time.Sleep(animationSpeed)
	}
}
