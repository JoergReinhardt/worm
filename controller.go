package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

func run() {

	// Allocate new game, pass termbox size as functional argument
	g := NewGame(termbox.Size)
	// render all worm segments and cherry once initially
	render(g)

	// game loop progresses the game in the background
	go func() {
		// set status to run
		g.state.stat = RUN
		for {
			// PAUSE MODE
			// if p is pressed, toggle game state and hold loop
			if g.state.stat == PAUSE {
				for {
					// check once per render cycle
					time.Sleep(animationSpeed)
					if g.state.stat != PAUSE {
						break
					}
				}
			}
			// PLAY
			//- detects colissions for next step
			//- ends the game on colission of worm with itself (sets
			//  state.stat to GAME_OVER)
			//- grows the worm and raises its speed on colission with cherry.
			//- moves the worm one step
			g.play()
			//- accesses game state to read cherry position, passes
			//callback closure to worms render method.
			render(g)
			// wait one worm speed duration until next step
			time.Sleep(g.state.speed)
		}
	}()

	// the event queue yields keyboard events
	eventQueue := make(chan termbox.Event)
	// dplete event queue in the background
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	// event retrieval and render loop. Blocks until game is stopped.
	for {
		ev := <-eventQueue
		if ev.Type == termbox.EventKey {
			switch {
			// set direction for the next move
			case ev.Key == termbox.KeyArrowUp || ev.Ch == 'k':
				g.state.move = UP
			case ev.Key == termbox.KeyArrowDown || ev.Ch == 'j':
				g.state.move = DOWN
			case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'h':
				g.state.move = LEFT
			case ev.Key == termbox.KeyArrowRight || ev.Ch == 'l':
				g.state.move = RIGHT
			case ev.Ch == 'p':
				// toggle game state between pause & run
				if g.state.stat == PAUSE {
					(*g.state).stat = RUN
				} else {
					(*g.state).stat = PAUSE
				}
			case ev.Key == termbox.KeyEsc || ev.Ch == 'q':
				return
			}
		}
		// exit, if game got ended by last move
		if g.state.stat == GAME_OVER {
			return
		}

		// render current game state
		render(g)
		// sleep til next animation cycle
		time.Sleep(animationSpeed)
	}
}
