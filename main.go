package main

import (
	"go-menu/router"
)

func main() {
	s := router.NewServer()
	s.Run(":8080")
}
