package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wearkinetic/awss3"
	"github.com/wearkinetic/keen"

	"github.com/wearkinetic/beutils"
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
	k, err := keen.NewFromEnv()
	company := keen.Company{
		Keen:       k,          // Keen instance
		Name:       os.Args[1], // company name as stored in keen
		ShiftHours: 8}          // shift length in hours
	if err != nil {
		log.Fatal(err)
	}

	client, errc := beutils.NewHTTPClient()
	if errc != nil {
		log.Fatal(errc)
	}

	devicesAtlocation, err := beutils.GetAllDevicesAtLocation(client, company.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("devicesAtlocation = ", devicesAtlocation)

	startDate := os.Args[2]
	endDate := os.Args[3]

	dates, err := beutils.GetDateRange(startDate, endDate)
	if err != nil {
		log.Fatal(err)
	}

	for _, checkdate := range dates {

		fmt.Printf("Checking for the date: %s\n", checkdate)

		response, err := company.GetData(
			checkdate+"T00:00:00-00:00", // start of timeframe to get
			checkdate+"T23:59:59-00:00", // end of timeframe to get
			"daily") // interval to group into
		if err != nil {
			log.Fatal(err)
		}

		// responseByTimeframe := response.ByTimeframeByEmployee() // group first by timeframe then by employee for easy marshalling to JSON
		// responseByEmployee := response.ByTimeframeByEmployee() // the reverse

		// fmt.Println(goutil.Pretty(*responseByEmployee))

		// --------
		// STEP 2. Setup S3 session
		// --------
		session := awss3.NewSession(awss3.REGION_US_EAST_1)

		// --------
		// STEP 4. For each device key assigned to the company, compare the data between Keen and S3
		// --------
		countNoData := 0
		for _, device := range devicesAtlocation {

			// First check if this device is assigned to an employee
			employeeName, employeeID := beutils.GetEmployeeInfo(client, device)
			if employeeName == "" {
				continue
			}

			lifts := 0
			activetime := 0

			for _, dt := range *response.Employees {
				if dt.ID == employeeID {
					lifts = dt.Lifts
					activetime = dt.ActiveSeconds
				}
			}

			list, err := session.List("kinetic-device-data", "raw/"+device+"/"+checkdate)
			if err != nil {
				log.Println("Couldn't read file list")
			}

			// Each file has either 1 data point (40ms) or 5mins (300s) of data
			// len(list)*0.04 <= activetime <= len(list)*5*60

			marker := ""
			if activetime > (len(list)*5*60 + 3600) { // Because there can be 1 hr of overlap from the other day in Keen data
				marker = "<--- missing S3 data"
			}

			if activetime < len(list)*1*6 && len(list) > 1 { // At least 1 sec of data and > 1 file in S3
				marker = "<--- missing Keen data"
			}

			if activetime > 3600*15 {
				marker = "<--- GREATER than 15 HRS"
			}

			if marker != "" {
				employeeName, _ := beutils.GetEmployeeInfo(client, device)
				fmt.Printf("%s\t%s\t%d\t%d\t%d\t%s\n", employeeName, device, lifts, activetime, len(list), marker)
			}

			if activetime == 0 && lifts == 0 && len(list) == 0 {
				countNoData++
			}
		}
		fmt.Printf("%d/%d with no data\n\n", countNoData, len(devicesAtlocation))
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
