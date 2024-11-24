package main

import (
	"checkmate/memory"
	"checkmate/optimization"
	"checkmate/routes"
)

func main() {

	// f, err := os.Create("cpu.prof")
	// if err != nil {
	// 	panic(err)
	// }
	// pprof.StartCPUProfile(f)

	// go func() {
	// 	time.Sleep(15 * time.Second)

	// 	pprof.StopCPUProfile()
	// 	fmt.Println("CPU profile stopped")
	// }()

	// database.Init()
	memory.Init()
	optimization.Init()
	routes.Serve()
}
