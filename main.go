package main

import (
	"go-menu/router"
	"log"
)

func main() {
	s := router.NewServer()
	if err := s.Run(":8080"); err != nil {
		log.Fatal("サーバーの起動に失敗しました: ", err)
	}
}
