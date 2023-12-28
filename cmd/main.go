package main

import (
	"dailybux/config"
	"dailybux/router"
)

func main() {
	config.Init()
	router.Make()
}
