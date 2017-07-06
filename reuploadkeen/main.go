package main

import (
	"fmt"

	"github.com/wearkinetic/logging"
	"github.com/wearkinetic/uploader/uploaderservice"
)

func main() {
	logging.Log(logging.INFO, fmt.Sprintf("Started uploader"))

	client, errc := uploaderservice.NewHTTPClient()
	if errc != nil {
		fmt.Printf("ERROR: %s\n", errc)
	}

	key := ""
	activityPath := "/Users/adityabansal/Downloads/home 5/kinetic/keen/2017-06-01T03:40:15.719Z"
	dirKeenPath := "/tmp"

	nfiles, errp := uploaderservice.ProcessActivityFiles(key, activityPath, dirKeenPath, client)
	if errp != nil {
		fmt.Printf("ERROR: %s\n", errp)
	}
	fmt.Printf("Uploaded %d events", nfiles)

}
