package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

//define structure for sec api response
type Company struct {
	Name         string
	Ticker       string
	CIK          string
	CUSIP        string
	Exchange     string
	IsDelisted   bool
	Category     string
	Sector       string
	Industry     string
	Sic          string
	SicSector    string
	SicIndustry  string
	FamaSector   string
	FamaIndustry string
	Currency     string
	Location     string
	Id           string
}

//get the SEC_API_TOKEN from .env file
func goDotEnvVariable(key string) string {

	//load .env file
	err := godotenv.Load(".env")

	//checks for existence and use of .env file
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//return the value of the key in question
	return os.Getenv(key)
}

//check if userNameInput is between 4 and 25 characters
func userNameInputLengthChecker(userNameInput string) bool {
	//set length of companySelection to variable companySelection
	var companySelectionLength = len(userNameInput)

	if companySelectionLength < 4 {
		fmt.Print("Company name input must be 4 or more characters long.")
	} else if companySelectionLength > 25 {
		fmt.Print("Company name input must be 25 or fewer characters long.")
	}
	return true
}

/*
//save this for later
//checks length of sec api response
func oneToManyHandler(responseLength int) {
	if responseLength == 0 {
		fmt.Print("\nNo company matches that description.\n")
	} else if responseLength == 1 {
		fmt.Print("\nNice! Found only one company matching that description:\n")
		//need to make a new function that checks length of ticker
		//tickerLengthChecker()
	} else {
		fmt.Printf("\nLooks like there are a few results, we'll need to narrow this down:\n")
		fmt.Println("\nPlease select one of the following options:")
		//tickerLengthChecker()
	}
	return
}
*/

//check length of ticker
func tickerLengthChecker(ticker string) bool {
	//if ticker length is greater than 0 and not N/A
	if len(ticker) > 0 && ticker != "N/A" {
		return true
	}
	return true
}

//calls sec api, returns ticker symbol and location of company
func mappingApiNameCaller(userNameInput string) string {

	if userNameInputLengthChecker(userNameInput) {

		fmt.Print("\nOkay, you want to know about ", userNameInput, "\n")

		//get the SEC_API_TOKEN from .env file to use in api call
		dotenv := goDotEnvVariable("SEC_API_TOKEN")

		//api call
		requestURL := fmt.Sprintf("https://api.sec-api.io/mapping/name/%s/?token=%s", userNameInput, dotenv) //userNameInput still the input
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}
		res, err := http.DefaultClient.Do(req)

		//check api response and present any errors
		if err != nil {
			log.Fatal(err)
		} else if res.StatusCode != 200 {
			fmt.Printf("Problem with the API request (status code: %d", res.StatusCode)
		}
		defer res.Body.Close()

		//read the api response and present any errors
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		//set variable to iterate through response
		var c []Company

		//unmarshal the json response of the body
		unmarshaledBody := json.Unmarshal(body, &c)
		if unmarshaledBody != nil {
			fmt.Println(unmarshaledBody)
		}

		var responseLength = len(c)

		///////////////////////////////////////////////////////////////////////////////////////////////////
		//add logic to execute different behavior for 1 vs many results
		//if number of results is 0, stop the loop

		if responseLength == 0 {
			fmt.Print("\nNo company matches that description.\n")
		} else if responseLength == 1 {
			fmt.Print("\nNice! Found only one company matching that description:\n")
			//no need to iterate through a single result here
			//probably cut this down to a simpler process
			for _, values := range c {
				//var ticker = values.Ticker
				if !tickerLengthChecker(values.Ticker) {
					fmt.Printf("\nLooks like %#v is not a public company (yet). Sorry!", userNameInput)
				} else {
					fmt.Printf("\n\tFull company name: %#v\n\tTicker: %#v\n\tLocation: %#v\n", values.Name, values.Ticker, values.Location)
				}
			}
			//if number of results is more than 1, break the loop
			//focus here
		} else {
			fmt.Printf("\nLooks like there are %d results, we'll need to narrow this down:\n", len(c))
			//fmt.Println("\nPlease select one of the following options:")
			//var userSelection int
			//fmt.Scanf("%g", &userSelection)
			for _, values := range c {
				if !tickerLengthChecker(values.Ticker) {

				} else {
					fmt.Printf("\n\tFull company name: %#v\n\tTicker: %#v\n\tLocation: %#v\n", values.Name, values.Ticker, values.Location)
				}
			}
		}
	} else {
		userNameInputLengthChecker(userNameInput)
	}

	//probably need to return the ticker and cik for future use...
	return ""
}

//Accepts user input about company,checks user input for length
func main() {

	//request user input about company
	fmt.Println("Which company do you want to know about?")

	//set user selected company as string variable
	var companySelection string
	fmt.Scanln(&companySelection)

	mappingApiNameCaller(companySelection)
}
