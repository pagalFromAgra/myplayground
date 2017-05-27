package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {

	log.Println(lastVersion())
}

// LastVersion
func lastVersion() string {
	vfile, err := os.Open("/Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/abplay/lastVersion/version")
	if err != nil {
		return "---"
	}
	defer vfile.Close()

	scanner := bufio.NewScanner(vfile)
	lastline := ""
	for scanner.Scan() {
		// Keep reading till we get to the last line
		lastline = fmt.Sprintf("%s", scanner.Text())
	}
	return lastline
}
