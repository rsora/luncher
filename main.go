package main

import (
	"log"
)

func main() {
	a := App{}
	a.Initialize()
	log.Println("Starting...")
	a.Run(":8000")
}
