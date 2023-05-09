package main

import (
	"cafe/config"
	"cafe/routes"
)

func main() {
	config.InitDB()
	e := routes.New()
	e.Start(":8000")
}
