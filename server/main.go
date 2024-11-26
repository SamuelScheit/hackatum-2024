package main

import (
	"checkmate/memory"
	"checkmate/optimization"
	"checkmate/routes"
	"os"
	"os/signal"
	"runtime/pprof"
)

func main() {

	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()

	// database.Init()
	memory.Init()
	optimization.Init()
	routes.Serve()
}
