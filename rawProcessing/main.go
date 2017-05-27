package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/wearkinetic/uploader/uploaderservice"
)

const (
	DATAFILES = "/Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/abplay/rawProcessing"
)

func main() {
	files, err := ioutil.ReadDir(DATAFILES)
	if err != nil {
		log.Println(err)
	}
	unprocessedCount := 0
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), "out") || file.IsDir() {
			continue
		}
		if file.Size() > 0 { // Sometimes it can be zero bytes long

			startMS, endMS, err := getStartEndTimes(file.Name())
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("[%s] start: %d & end: %d\n", file.Name(), startMS, endMS)

			unprocessedCount++
		}
		// os.Rename(path.Join(DATAFILES, file.Name()), path.Join(DATAFILES, file.Name()+".processed"))
	}

}

func getStartEndTimes(file string) (uint64, uint64, error) {
	// Read the start and end times
	fd, err := os.Open(file)
	if err != nil {
		return 0, 0, err
	}
	return uploaderservice.ExtractStartEndTimes(fd)
}
