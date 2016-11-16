package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/nsf/termbox-go"
)

const animationSpeed = 10 * time.Millisecond

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	err := termbox.Init()
	//_ = termbox.Output256
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	run()
}
