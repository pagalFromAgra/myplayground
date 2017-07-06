package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	DEFAULT_API_URL = "http://beasag-develop.wearkinetic.com"

	wuStatusURL = DEFAULT_API_URL + "/v1/wus"

	reportURL          = DEFAULT_API_URL + "/v1/wilson"
	pendingCommandsURL = DEFAULT_API_URL + "/v1/pendingCommands"
	bootstrapURL       = DEFAULT_API_URL + "/v1/wus"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Need location to papertrail output")
		os.Exit(0)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	client, errc := NewHTTPClient()
	if errc != nil {
		log.Fatal(errc)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		valarray := strings.Fields(scanner.Text())

		// Remove "-1" or "-2" etc. from the key added by papertrail
		key := strings.Split(valarray[1], "-")

		location, sku := getLocationSKU(client, key[0])
		fmt.Printf("2017-06-%s,%s,%s,%s,%s\n", os.Args[2], key[0], valarray[0], sku, location)
	}

}

func getLocationSKU(client *http.Client, device string) (string, string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?uid=%s", wuStatusURL, device), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoic3VwZXJ1c2VyIn0.VHLJ1X_FvBq_Lu275k4WoJB2LD4LTg49cuxbj6c5Ru0")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(responseData))

	var dat map[string]interface{}
	var wuarray []interface{}
	var wudata map[string]interface{}

	err = json.Unmarshal(responseData, &dat)
	if err != nil {
		log.Fatal(err)
	}

	wuarray = dat["data"].([]interface{})

	wudata = wuarray[0].(map[string]interface{})

	return wudata["locationUID"].(string), wudata["sku"].(string)
}
