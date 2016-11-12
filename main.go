package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

const animationSpeed = 10 * time.Millisecond

var wormSpeed = 250 * time.Millisecond

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

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
			time.Sleep(wormSpeed)
			// exit, if game over
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
