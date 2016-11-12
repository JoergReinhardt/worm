package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

func run() {

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	g := NewGame()
	render(g)

	// game loop (background)
	go func() {
		(*g.state).stat = RUN
		for {
			(*g).move()
			render(g)
			// move with worm speed
			time.Sleep((*g.state).speed)

			// if p is pressed, hold worm
			if (*g.state).stat == PAUSE {
				for {
					time.Sleep(animationSpeed)
					if (*g.state).stat != PAUSE {
						break
					}
				}
			}
		}
	}()

	// event/render loop (blocks)
	for {
		ev := <-eventQueue
		if ev.Type == termbox.EventKey {
			switch {
			case ev.Key == termbox.KeyArrowUp || ev.Ch == 'k':
				(*g.state).move = UP
			case ev.Key == termbox.KeyArrowDown || ev.Ch == 'j':
				(*g.state).move = DOWN
			case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'h':
				(*g.state).move = LEFT
			case ev.Key == termbox.KeyArrowRight || ev.Ch == 'l':
				(*g.state).move = RIGHT
			case ev.Ch == 'p':
				// toggle pause/run
				if g.state.stat == PAUSE {
					(*g.state).stat = RUN
				} else {
					(*g.state).stat = PAUSE
				}
			case ev.Key == termbox.KeyEsc || ev.Ch == 'q':
				return
			}
		}
		// exit, if game over
		if (*g.state).stat == GAME_OVER {
			return
		}

		render(g)
		time.Sleep(animationSpeed)
	}
}
