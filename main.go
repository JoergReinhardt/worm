package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

const animationSpeed = 250 * time.Millisecond

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
				(*g.state).stat = PAUSE
			case ev.Key == termbox.KeyEsc || ev.Ch == 'q':
				return
			}
		}
		(*g).move()
		render(g)
		time.Sleep(animationSpeed)
	}
}
