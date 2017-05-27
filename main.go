package main

import (
	"log"
	"os"
)

func main() {
	dirpath := "/Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/abplay/testdir"
	os.RemoveAll(dirpath)
	log.Println("Removed ", dirpath)
}
