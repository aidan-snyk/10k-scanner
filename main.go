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

//calls sec api, returns ticker symbol and location of company
func mappingApiNameCaller(userNameInput string) string {

	//set length of companySelection to variable companySelection
	var companySelectionLength = len(userNameInput)

	//check if companySelection is between 4 and 25 characters
	if companySelectionLength >= 4 && companySelectionLength <= 25 {

		fmt.Print("\nOkay, you want to know about ", userNameInput, "\n")

		//get the SEC_API_TOKEN from .env file to use in api call
		dotenv := goDotEnvVariable("SEC_API_TOKEN")

		//alternate api call using http.NewRequest
		requestURL := fmt.Sprintf("https://api.sec-api.io/mapping/name/%s/?token=%s", userNameInput, dotenv)
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}
		res, err := http.DefaultClient.Do(req)

		fmt.Printf("client: got response!\n")
		fmt.Printf("client: status code: %d\n", res.StatusCode)

		//check api response and present any errors
		if err != nil {
			log.Fatal(err)
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

		//add logic to execute different behavior for 1 vs many results
		//if number of results is 0, stop the loop
		if len(c) == 0 {
			fmt.Print("\nHmmmm not finding a company that matches that description.\n")

			//if number of results is 1, present only that one
		} else if len(c) == 1 {
			fmt.Print("\nNice! Found only one company matching that description:\n")
			for _, values := range c {
				if len(values.Ticker) == 0 {
					fmt.Printf("\nLooks like %#v is not a public company (yet). Sorry!", userNameInput)
				} else {
					fmt.Printf("\nTicker: %#v\nLocation: %#v\n", values.Ticker, values.Location)
				}
			}

			//if number of results is more than 1, break the loop
		} else {
			fmt.Printf("\nLooks like there are %d results, we'll need to narrow this down:\n", len(c))
			for _, values := range c {
				if len(values.Ticker) == 0 || values.Ticker == "N/A" {
					continue
				} else {
					fmt.Printf("\nTicker: %#v\nLocation: %#v\n", values.Ticker, values.Location)
				}
			}
		}
		//reject user input that is not between 4 and 25 characters (inclusive)
	} else {
		fmt.Println("Your entry must be between 4 and 25 characters long, please try again")
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
