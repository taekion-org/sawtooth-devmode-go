package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	// initialize RNG
	rand.Seed(time.Now().UnixNano())
}

func main() {
	fmt.Println("Hello, Gopher!")
}
