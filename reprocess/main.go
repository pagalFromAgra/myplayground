package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	OUT_DIR     = "/Users/adityabansal/digger/catclayton/reprocess-all/catclayton-reprocess/catclayton_left_2017-04-27_2017-05-13_reprocess"
	TIME_FORMAT = "2006-01-02T15:04:05"
)

// Message structure, service agnostic
type Message struct {
	Action string `json:"action"`
	Body   struct {
		Start int64  `json:"start"`
		End   int64  `json:"end"`
		Type  string `json:"type"`
	} `json:"body"`
}

func main() {
	files, err := ioutil.ReadDir(OUT_DIR)
	if err != nil {
		log.Println("Could not read dir, error: ", err)
	}

	for _, file := range files {
		splitted := strings.Split(file.Name(), ".")
		deviceKey := splitted[0]

		filepath := fmt.Sprintf("%s/%s", OUT_DIR, file.Name())

		readConvert(filepath, deviceKey)
	}

}

// Parse a log file, and add service, then dump to other file
func readConvert(filepath string, deviceKey string) error {

	// We'll save all the data in a bytes buffer
	// body := &bytes.Buffer{}
	// writer := multipart.NewWriter(body)
	// part, err := writer.CreateFormFile("file", filepath)

	f, err := os.Open(filepath)
	if err != nil {
		log.Println("Could not open file", filepath)
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	msg := &Message{}
	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), msg)
		if err != nil {
			continue
		}

		if msg.Action == "lift.risky" && (msg.Body.Type == "twist" || msg.Body.Type == "bend") {

			// s, _ := time.Parse(TIME_FORMAT, strconv.FormatUint(msg.Body.Start, 10))
			// e, _ := time.Parse(TIME_FORMAT, strconv.FormatUint(msg.Body.End, 10))

			// if errs != nil || erre != nil {
			// 	continue
			// }

			// log.Println(s.String())

			fmt.Printf("%s, %s, %s, %s\n", deviceKey, msg.Body.Type, time.Unix(msg.Body.Start/1000, 0).Format(TIME_FORMAT), time.Unix(msg.Body.End/1000, 0).UTC().Format(TIME_FORMAT))
		}
	}

	return nil
}
