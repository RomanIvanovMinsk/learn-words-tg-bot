package elasticService

import (
	"bytes"     // for converting JSON to bytes array
	"fmt"       // for printing to console
	"io/ioutil" // for reading IO of JSON file
	"log"       // for logging errors
	"net/http"  // for making HTTP requests
	"os"        // for opening JSON file
	"reflect"   // get object type
)

func ParseJsonAndSendItInElastick() {

	// Allow for custom formatting of log output
	log.SetFlags(0)

	file := OpenFile()
	defer file.Close()
	data := ReadFile(file)
	req := setUpRequestForFile(data)
	InvokeReqReadResp(req)
}

func setUpRequestForFile(byteSlice []byte) *http.Request {
	// Make HTTP request using "PUT" or "POST" verb
	req, err := http.NewRequest("PUT", "http://localhost:9200/_bulk?pretty=true", bytes.NewBuffer(byteSlice))

	// ES 6.0> requires Content-Type header to avoid 406 HTTP error:
	// "error":"Content-Type header [] is not supported","status":406}
	req.Header.Set("Content-Type", "application/x-ndjson")

	// Print out the HTTP request and check for errors
	if err != nil {
		log.Fatalf("http.NewRequest ERROR:", err)
	} else {
		fmt.Println("HTTP Request:", req)
	}

	return req
}

func ReadFile(file *os.File) []byte {

	// Call ioutil.ReadAll() to create a bytes array from file's JSON data
	byteSlice, err := ioutil.ReadAll(file)

	// Check for IO errors
	if err != nil {
		log.Fatalf("ioutil.ReadAll() ERROR:", err)
	}
	fmt.Println("bytesStr TYPE:", reflect.TypeOf(byteSlice), "n")

	return byteSlice
}

func OpenFile() *os.File {
	// Use the OS package to load the JSON file
	file, err := os.Open("docs.json")
	if err != nil {
		log.Fatalf("os.Open() ERROR:", err)
	}
	// Close the file AFTER operations are complete
	return file
}

func InvokeReqReadResp(req *http.Request) {
	// Instantiate a new client object
	client := &http.Client{}

	// Pass HTTP request to Elasticsearch client and check for errors
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.DoERROR:", err)
	}

	// Close the response body after operations are complete
	defer resp.Body.Close()

	// Parse out the response body and check for errors
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("client.Do ERROR:", err)
	}

	// Convert the bytes object []uint8 of the JSON response to a string
	strBody := string(body)

	// Print out the response body
	fmt.Println("nresp.Body:", strBody)
}
