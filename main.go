package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
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

type tenKayFiling struct {
	Total struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
	Query struct {
		From int `json:"from"`
		Size int `json:"size"`
	} `json:"query"`
	Filings []struct {
		ID                  string `json:"id"`
		AccessionNo         string `json:"accessionNo"`
		Cik                 string `json:"cik"`
		Ticker              string `json:"ticker"`
		CompanyName         string `json:"companyName"`
		CompanyNameLong     string `json:"companyNameLong"`
		FormType            string `json:"formType"`
		Description         string `json:"description"`
		FiledAt             string `json:"filedAt"`
		LinkToTxt           string `json:"linkToTxt"`
		LinkToHTML          string `json:"linkToHtml"`
		LinkToXbrl          string `json:"linkToXbrl"`
		LinkToFilingDetails string `json:"linkToFilingDetails"`
		Entities            []struct {
			CompanyName          string `json:"companyName"`
			Cik                  string `json:"cik"`
			IrsNo                string `json:"irsNo"`
			StateOfIncorporation string `json:"stateOfIncorporation"`
			FiscalYearEnd        string `json:"fiscalYearEnd"`
			Type                 string `json:"type"`
			Act                  string `json:"act"`
			FileNo               string `json:"fileNo"`
			FilmNo               string `json:"filmNo"`
			Sic                  string `json:"sic"`
		} `json:"entities"`
		DocumentFormatFiles []struct {
			Sequence    string `json:"sequence"`
			Description string `json:"description,omitempty"`
			DocumentURL string `json:"documentUrl"`
			Type        string `json:"type"`
			Size        string `json:"size"`
		} `json:"documentFormatFiles"`
		DataFiles []struct {
			Sequence    string `json:"sequence"`
			Description string `json:"description"`
			DocumentURL string `json:"documentUrl"`
			Type        string `json:"type"`
			Size        string `json:"size"`
		} `json:"dataFiles"`
		SeriesAndClassesContractsInformation []interface{} `json:"seriesAndClassesContractsInformation"`
		PeriodOfReport                       string        `json:"periodOfReport"`
	} `json:"filings"`
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

func userNameInputLengthChecker(userNameInput string) bool {
	var companySelectionLength = len(userNameInput)

	if companySelectionLength < 4 || companySelectionLength > 25 {
		return false
	}
	return true
}

//check length and content of ticker
func tickerLengthChecker(ticker string) bool {
	//ensure ticker is greater than 0 characters and not "N/A"
	if len(ticker) == 0 || ticker == "N/A" {
		return false
	}
	return true
}

//calls sec api, returns ticker symbol and location of company
func mappingApiNameCaller(userNameInput string) string {

	//if userNameInput is the correct length, move forward
	if userNameInputLengthChecker(userNameInput) {

		fmt.Print("\nOkay, you want to know about ", userNameInput, "\n")
		//get the SEC_API_TOKEN from .env file to use in api call
		dotenv := goDotEnvVariable("SEC_API_TOKEN")

		//api call
		requestURL := fmt.Sprintf("https://api.sec-api.io/mapping/name/%s/?token=%s", userNameInput, dotenv)
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

		//add logic to execute different behavior for 1 vs many results
		//if number of results is 0, stop the loop
		if responseLength == 0 {
			fmt.Print("\nNo company matches that description.\n")

		} else if responseLength == 1 {
			//no need to iterate through a single result here
			//probably cut this down to a simpler process
			for _, values := range c {
				//var ticker = values.Ticker
				if tickerLengthChecker(values.Ticker) {
					fmt.Print("\nNice! Found only one company matching that description:\n")
					fmt.Printf("\n\tFull company name: %#v (%#v)\n\tTicker: %#v\n\tLocation: %#v\n",
						values.Name,
						values.CIK,
						values.Ticker,
						values.Location)
					return values.CIK
				} else {
					fmt.Printf("\nLooks like %#v is not a public company (yet). Sorry!", userNameInput)
					return ""
				}
			}
			//if number of results is more than 1, make user select
		} else {
			fmt.Print("\nA few companies match that name:\n")
			//set s as slice to hold tickers
			var s []string

			//iterate through all companies matching user input
			for _, values := range c {

				//skip if ticker doesn't meet public company conditions
				if tickerLengthChecker(values.Ticker) {
					//append company ticker symbol to s
					s = append(s, values.CIK)

					//present option number, name, ticker, and location to user
					fmt.Printf("\nOption %#v\n\tFull company name: %#v (CIK: %#v,Ticker: %#v)\n\tLocation: %#v\n",
						len(s),
						values.Name,
						values.CIK,
						values.Ticker,
						values.Location)
				}
			}
			var i int
			fmt.Println("\nEnter the number option you'd like to see: ")
			fmt.Scanln(&i)
			return s[i-1]
		}
	} else {
		fmt.Print("Company name input must be between 4 and 25 characters long.")
		return ""
	}
	return ""
}

/*
//gets link to filing text
func getFilingTxt(tickerSelection string) string {
	dotenv := goDotEnvVariable("SEC_API_TOKEN")


}
*/

//Accepts user input about company,checks user input for length
func main() {

	//request user input about company
	fmt.Println("Which company do you want to know about (no spaces please)?")

	//set user selected company as string variable
	var companySelection string
	fmt.Scanln(&companySelection)

	//testing adding the returned ticker to jsonData
	tickerSelection := mappingApiNameCaller(companySelection)

	//adding filing text caller
	dotenv := goDotEnvVariable("SEC_API_TOKEN")
	httpposturl := fmt.Sprintf("https://api.sec-api.io?token=%s", dotenv)

	var jsonData = []byte(fmt.Sprintf(`{
		"query": {
			"query_string": {
				"query": "cik:\"%s\" AND formType:\"10-K\""
			}
		},
		"from": "0",
		"size": "1",
		"sort": [{ "filedAt": { "order": "desc" } }]
	  }`, tickerSelection))

	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer((jsonData)))
	if error != nil {
		fmt.Printf("client: could not create request %s\n", error)
		fmt.Println(string(jsonData))
		os.Exit(1)
	}
	request.Header.Set("Content-Type", "application/json; chartset=UTF-8")

	client := &http.Client{}

	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	//return the link to the latest 10-K filing
	result := gjson.Get(string(body), "filings.#.linkToTxt")

	for _, filing := range result.Array() {
		println(filing.String())
	}
}
