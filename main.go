package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

const animationSpeed = 10 * time.Millisecond

func main() {
	err := termbox.Init()
	//_ = termbox.Output256
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	run()
}
