package elasticService

import (
	model "WordsBot/models"
	"bytes" // for converting JSON to bytes array
	"encoding/json"
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
	//DropOldIndex()
	//SetUpIndex()
	UploadDataToIndex()

}

// func DropOldIndex(){

// }

func SetUpIndex() {
	settings := &model.SettingsModel{}
	settings.Settings.NumberOfShards = 1

	req := SetUpRequestForCreateIndex(settings)
	InvokeReqReadResp(req)
}

func UploadDataToIndex() {
	file := OpenFile()
	defer file.Close()
	data := ReadFile(file)
	req := setUpRequestForFile(data)
	InvokeReqReadResp(req)
}

func SetUpRequestForCreateIndex(payload interface{}) *http.Request {
	reqBytes, err := json.Marshal(payload)

	//map from one json to another
	// mapReqBytes

	// Make HTTP request using "PUT" or "POST" verb
	req, err := http.NewRequest("PUT", "http://localhost:9200/dict", bytes.NewBuffer(reqBytes))
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

func setUpRequestForFile(byteSlice []byte) *http.Request {
	var result map[string]string
	json.Unmarshal([]byte(byteSlice), &result)
	indexedWords := make([]byte, 1024)
	i := 0
	for key, value := range result {
		index := &model.IndexedModel{}
		index.Index.Id = i
		index.Index.Index = "dict"
		result, err := json.Marshal(index)
		if err != nil {
			continue
		}

		word := &model.DictionaryWord{}
		word.Key = key
		word.Value = value

		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(word)

		indexedWords = append(indexedWords, result...)
		//indexedWords = append(indexedWords, []byte("\n")...)
		indexedWords = append(indexedWords, reqBodyBytes.Bytes()...)
		//indexedWords = append(indexedWords, []byte("\n")...)
	}

	// Make HTTP request using "PUT" or "POST" verb
	req, err := http.NewRequest("PUT", "http://localhost:9200/_bulk?pretty=true", bytes.NewBuffer(indexedWords))

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
	file, err := os.Open("./dict/dictionary.json")
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
