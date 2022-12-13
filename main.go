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

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

/*
	Accepts user input about a desired company,
	checks user input for length,
	calls sec api to get information about company,
	returns ticker symbol and location of company
*/
func main() {

	//request user input about company
	fmt.Println("Which company do you want to know about?")

	//set user selected company as string variable
	var companySelection string
	fmt.Scanln(&companySelection)

	//set length of companySelection to variable companySelection
	var companySelectionLength = len(companySelection)

	//check if companySelection is between 4 and 25 characters
	if companySelectionLength >= 4 && companySelectionLength <= 25 {
		fmt.Print("\nOkay, you want to know about ", companySelection, "\n")

		//get the SEC_API_TOKEN from .env file to use in api call
		dotenv := goDotEnvVariable("SEC_API_TOKEN")

		//call the sec api to get information about companySelection
		res, err := http.Get("https://api.sec-api.io/mapping/name/" + companySelection + "?token=" + dotenv)

		//handle the errors
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var c []Company

		unmarshaledBody := json.Unmarshal(body, &c)
		if unmarshaledBody != nil {
			fmt.Println(unmarshaledBody)
			return
		}

		for _, values := range c {
			fmt.Printf("\nTicker: %#v\nLocation: %#v\n", values.Ticker, values.Location)
		}
	} else {
		fmt.Println("Your entry must be between 4 and 25 characters long, please try again")
	}
}
