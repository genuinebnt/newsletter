package main

import (
	lib "genuinebnt/newsletter/internal"
)

func main() {
	server := lib.Server()
	server.Run(":8000")
}
