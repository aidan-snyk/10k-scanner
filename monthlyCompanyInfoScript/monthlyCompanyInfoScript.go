package main

import (
	"io"
	"net/http"
	"os"
)

//checks input for errors
func check(e error) {
	if e != nil {
		//code prints a panic message and function crashes
		panic(e)
	}
}

func main() {
	//SEC API endpoint for company information
	resp, err := http.Get("https://www.sec.gov/files/company_tickers.json")
	check(err)

	defer resp.Body.Close()

	//create a new file to hold results
	out, err := os.Create("test.json")
	check(err)

	defer out.Close()

	//write JSON results to test.json
	io.Copy(out, resp.Body)
}
