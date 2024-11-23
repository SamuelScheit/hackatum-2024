package main

import (
	"checkmate/database"
	"checkmate/optimization"
	"checkmate/routes"
)

func main() {
	database.Init()
	optimization.Init()
	routes.Serve()
}
