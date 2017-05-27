package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/wearkinetic/awss3"
	"github.com/wearkinetic/keendevice"
)

type Result struct {
	ID        string `json:"id"`
	Timeframe struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}
	Lifts       int `json:"lifts"`
	Time_Active int `json:"time_active"`
	Lift_Rate   int `json:"lift_rate"`
}

func main() {

	// --------
	// STEP 1. Get the data from Keen
	// --------
	k, _ := keendevice.NewFromEnv()
	company := keendevice.Company{
		Keen:       k,          // Keen instance
		Name:       "irmshred", // company name as stored in keen
		ShiftHours: 8}          // shift length in hours

	checkdate := "2017-05-26"

	response, _ := company.GetData(
		checkdate+"T00:00:00-00:00", // start of timeframe to get
		checkdate+"T23:59:59-00:00", // end of timeframe to get
		"daily") // interval to group into

	// responseByTimeframe := response.ByTimeframeByEmployee() // group first by timeframe then by employee for easy marshalling to JSON
	// responseByEmployee := response.ByTimeframeByEmployee() // the reverse

	// fmt.Println(goutil.Pretty(*responseByEmployee))

	// --------
	// STEP 2. Setup S3 session
	// --------
	session := awss3.NewSession(awss3.REGION_US_EAST_1)

	// --------
	// STEP 3. Read the device keys from the exported employee csv
	// --------
	f, _ := os.Open("/Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/digger/export-irm.csv")
	r := csv.NewReader(bufio.NewReader(f))
	result, _ := r.ReadAll()

	// --------
	// STEP 4. For each device key assigned to the company, compare the data between Keen and S3
	// --------
	for _, row := range result {

		device := row[2]
		lifts := 0
		activetime := 0

		for _, dt := range *response.Employees {
			if dt.ID == device {
				lifts = dt.Lifts
				activetime = dt.ActiveSeconds
			}
		}

		list, err := session.List("kinetic-device-data", "raw/"+device+"/"+checkdate)
		if err != nil {
			log.Println("Couldn't read file list")
		}

		marker := ""
		if activetime > len(list)*5*60 {
			marker = "<--- missing S3 data"
		}

		if activetime < len(list)*5*48 {
			marker = "<--- missing Keen data"
		}

		fmt.Printf("%s\t%d\t%d\t%d\t%s\n", device, lifts, activetime, len(list), marker)
	}
	//
	// for _, row := range allinfo {
	//
	// 	if strings.Contains(row.side, "left") {
	// 		fmt.Printf("%s;", row.device)
	// 	}
	// }

	// i := 0
	// for _, dt := range *response.Employees {
	// 	// fmt.Printf("%d %s %d %d\n", i, dt.ID, dt.Lifts, dt.ActiveSeconds)
	// 	fmt.Println(dt)
	// 	i++
	// }

}
