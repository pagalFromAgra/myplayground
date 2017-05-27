package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/ironbay/dynamic"
	"github.com/wearkinetic/drs/drs-go"
	"github.com/wearkinetic/go/core/domains/device"
	"github.com/wearkinetic/go/core/domains/ipc"
	"github.com/wearkinetic/logging"
)

// BLACKLIST is the list of all the events which are high frequency and should be ignored
var BLACKLIST = []string{
	"i2c",
}

const KEEN_LAYOUT = "2006-01-02T15:04:05.000Z"

// msInHr is the number of milliseconds in an hour
var msInHr int64 = 3600000

func main() {

	// if err := os.MkdirAll(device.PATH_ACTIVITY, 0777); err != nil {
	// 	log.Fatal(err)
	// }
	//
	// satellite := ipc.New()
	//
	// company, employee, job := readConfig()
	//
	// // Handshake with uart. Don't record events if feedback is muted
	// logEvents := true // by default log the events
	// satellite.On("feedback.mute", func(cmd *drs.Command) {
	// 	logging.Log(logging.INFO, fmt.Sprintln("Received feedback.mute"))
	// 	if strings.Contains(cmd.Body.(string), "true") {
	// 		logEvents = false
	// 	} else {
	// 		logEvents = true
	// 	}
	// 	sendAck(satellite, logEvents)
	// })
	//
	// // Always read the configuration file when going off the dock
	// // to get the latest company, employee and job info
	// satellite.On("dock.off", func(cmd *drs.Command) {
	// 	company, employee, job = readConfig()
	// })
	//
	// // On the event "lift.risky.raw"
	// // Device must be assigned to an employee
	// satellite.On("lift.risky", func(cmd *drs.Command) {
	// 	// If employee name doesn't exist, read the config file
	// 	if logEvents {
	// 		if employee == "" {
	// 			company, employee, job = readConfig()
	// 		}
	// 		if employee != "" {
	// 			recordRawActivity("lift.risky", cmd, company, employee, job)
	// 		}
	// 	}
	// })
	//
	// // On the event "active.interval"
	// // Device must be assigned to an employee
	// satellite.On("active.interval", func(cmd *drs.Command) {
	// 	if logEvents {
	// 		if employee == "" {
	// 			company, employee, job = readConfig()
	// 		}
	// 		if employee != "" {
	// 			if err := sliceActiveTimeHourly("active.interval", cmd, company, employee, job); err != nil {
	// 				logging.Log(logging.ERROR, fmt.Sprintln("in sliceActiveTimeHourly: ", err))
	// 			}
	// 		}
	// 	}
	// })
	//
	// satellite.Read()

	printIntervals(1494331328455, 1494333754725)

}

func tFromF(f int64) time.Time {
	tInt := int64(f)
	millisec := tInt % 1000
	sec := (tInt - millisec) / 1000
	nsec := millisec * 1000000
	return time.Unix(sec, nsec).UTC()
}

func printIntervals(start, end int64) error {

	// write the duration for easier calculation
	value := end - start
	if value < 0 {
		return errors.New("ERROR: end < start")
	}

	subIntervalStart := start
	subIntervalEnd := nextClockHourEnd(subIntervalStart)

	// For cases within an hour
	if end <= subIntervalEnd {
		subIntervalEnd = end
	}

	for subIntervalEnd <= end {

		fmt.Printf("Interval: [%s, %s]\n", tFromF(subIntervalStart).UTC().Format(KEEN_LAYOUT), tFromF(subIntervalEnd).UTC().Format(KEEN_LAYOUT))

		if subIntervalEnd == end {
			break
		} // break out if the last interval

		// Setup the next interval
		subIntervalStart = subIntervalEnd
		if (subIntervalEnd + msInHr) > end { // i.e. last interval
			subIntervalEnd = end
		} else {
			subIntervalEnd += msInHr
		}
	}

	return nil
}

func sendAck(satellite *ipc.IPC, logEvents bool) {
	if logEvents {
		satellite.Fire(&drs.Command{
			Action: "log.events",
			Body:   "on",
		})
	} else {
		satellite.Fire(&drs.Command{
			Action: "log.events",
			Body:   "off",
		})
	}
}

func readConfig() (string, string, string) {
	config, err := device.Config()
	if err != nil {
		return "", "", ""
	}
	return config["company"], config["employee"], config["job"]
}

func recordRawActivity(activityName string, cmd *drs.Command, company string, employee string, job string) {
	data := cmd.Map()
	start := dynamic.Int(data, "start")
	end := dynamic.Int(data, "end")
	value := dynamic.Int(data, "count")
	if value == 0 {
		value = 1
	}
	riskType := dynamic.String(data, "type")
	if riskType == "" {
		riskType = "na"
	}
	fileName := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v|%v", activityName, value, start, end, company, employee, job, riskType)
	os.Create(path.Join(device.PATH_ACTIVITY, fileName))
}

func sliceActiveTimeHourly(intervalName string, cmd *drs.Command, company string, employee string, job string) error {
	data := cmd.Map()
	start := dynamic.Int(data, "start")
	end := dynamic.Int(data, "end")

	// write the duration for easier calculation
	value := end - start
	if value < 0 {
		return errors.New("ERROR: end < start")
	}

	subIntervalStart := start
	subIntervalEnd := nextClockHourEnd(subIntervalStart)

	// For cases within an hour
	if end <= subIntervalEnd {
		subIntervalEnd = end
	}

	for subIntervalEnd <= end {

		recordActiveInterval(intervalName, int64((subIntervalEnd-subIntervalStart)/1000), subIntervalStart, subIntervalEnd, company, employee, job)

		if subIntervalEnd == end {
			break
		} // break out if the last interval

		// Setup the next interval
		subIntervalStart = subIntervalEnd
		if (subIntervalEnd + msInHr) > end { // i.e. last interval
			subIntervalEnd = end
		} else {
			subIntervalEnd += msInHr
		}
	}

	return nil
}

func nextClockHourEnd(start int64) int64 {
	return start + (msInHr - (start % msInHr))
}

func recordActiveInterval(activity string, value int64, start int64, end int64, company string, employee string, job string) {

	// Some checks
	if value < 0 || value > 3600 {
		logging.Log(logging.ERROR, fmt.Sprintf("active interval outside bound (0<=x<=3600), %d ", value))
		// Keep moving and save this value in Keen just in case
	}

	fileName := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v|%v", activity, value, start, end, company, employee, job, "na") // na is the Info field
	os.Create(path.Join(device.PATH_ACTIVITY, fileName))
}
