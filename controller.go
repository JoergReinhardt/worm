package main

import (
	"github.com/nsf/termbox-go"
	"time"
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

type GameStat uint8

const (
	INIT GameStat = 0
	RUN  GameStat = 1 << iota
	PAUSE
	GAME_OVER
)

// the gameController runs worm and cherry at the rate required by current worm
// speed
func gameController(g *Game) {
	// set status to run
	// g.state.stat = RUN
	for {
		// INIT SCREEN
		if g.state.eventState == INIT {
			for { // wait for game start or quit
				initScreen()
				// check once per render cycle
				time.Sleep(animationSpeed)
				if g.state.eventState != INIT {
					break
				}
			}
		}
		// PAUSE MODE
		// if p is pressed, toggle game state and hold loop
		if g.state.eventState == PAUSE {
			for {
				initScreen()
				// check once per render cycle
				time.Sleep(animationSpeed)
				if g.state.eventState != PAUSE {
					break
				}
			}
		}
		// PLAY
		//- detects colissions for next step
		//- ends the game on colission of worm with itself (sets
		//  state.stat to GAME_OVER)
		//- grows the worm and raises its speed on colission
		//  with cherry.
		//- moves the worm one step
		g.play()
		//- accesses game state to read cherry position
		//- use termbox SetCell, to render cherry
		//- pass termbox SetCell to worms render method.
		render(g)
		// wait one worm speed duration cycle until next move (event
		// queue is running parallel meanwhile changing the game state.
		time.Sleep(g.state.speed)
	}
}

func run() { // runs the animation and input event cycles

	//////////////////////////////////////////////////////////////////////////////////////
	/// GAME SETUP
	// Allocate new game, pass termbox size as functional argument
	g := NewGame(termbox.Size)
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
			case ev.Key == termbox.KeyArrowUp || ev.Ch == 'k':
				g.state.direction = UP
			case ev.Key == termbox.KeyArrowDown || ev.Ch == 'j':
				g.state.direction = DOWN
			case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'h':
				g.state.direction = LEFT
			case ev.Key == termbox.KeyArrowRight || ev.Ch == 'l':
				g.state.direction = RIGHT
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
			case ev.Key == termbox.KeyEsc || ev.Ch == 'q':
				return
			}
		}
		// exit, if game got ended by last move
		if g.state.eventState == GAME_OVER {
			(*g).reset()
			(*g.state).eventState = INIT
		}

		// render current game state
		render(g)

		// sleep til next animation cycle
		time.Sleep(animationSpeed)
	}
}
