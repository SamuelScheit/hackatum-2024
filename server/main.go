package main

import (
	"checkmate/database"
	"checkmate/routes"
)

func main() {
	database.Init()
	routes.Serve()
}
