package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	i := 0
	// var list []string
	for {
		i++
		fmt.Println("now cnt:", i)
		log.Println("now cnt:", i)
		time.Sleep(time.Second * 2)
		// log.Println("list", list[2])
	}
	select {}
}
