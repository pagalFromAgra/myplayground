package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func main() {
	f, _ := os.Open("/Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/abplay/readcsv/employeelist.csv")

	r := csv.NewReader(bufio.NewReader(f))
	result, _ := r.ReadAll()

	for _, row := range result {
		if strings.Contains(row[9], "left") {
			log.Println(row[0])
		}
	}
}
