package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"os"

	"io/ioutil"

	"encoding/json"

	"encoding/csv"

	"fmt"
	"io"

	"golang.org/x/net/publicsuffix"
)

var clientProduction http.Client
var clientStage http.Client

var configParameters configJson

var rows [][]string

func main() {
	// Login to ruby
	log.Println("Start...")

	bootstrap()

	log.Println("Login to production...")
	if err := login(&clientProduction, 0); nil != err {
		log.Fatal(err)
	}
	defer logout(&clientProduction, 1)

	log.Println("Login to stage...")
	if err := login(&clientStage, 1); nil != err {
		log.Fatal(err)
	}
	defer logout(&clientStage, 1)

	log.Println("Creating results file...")
	resultsFile, err := os.Create(configParameters.Results)
	if err != nil {
		log.Fatalf("Error creating results file (err: %s)", err)
	}
	defer resultsFile.Close()

	log.Println("Reading endpoints file...")
	endpointsFile, err := os.Open(configParameters.Endpoints)
	if err != nil {
		log.Fatalf("Error creating results file (err: %s)", err)
	}
	defer endpointsFile.Close()

	reader := csv.NewReader(endpointsFile)

	for {
		endpoint, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		diffEndpoint(endpoint[0])
	}

	log.Println("End.")
}

func bootstrap() {
	configFile, err := ioutil.ReadFile(jsonFileName)
	if nil != err {
		log.Fatalf("Error opening config file: %v", err)
	}

	err = json.Unmarshal(configFile, &configParameters)
	if nil != err {
		log.Fatalf("Error at the unmarshal: %v", err)
	}
}

func login(c *http.Client, environment int) error {
	var err error

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)
	if nil == err {
		loginUrl := new(url.URL)
		loginUrl.Scheme = "http"
		loginUrl.Host = configParameters.Domains[environment]
		loginUrl.RawPath = configParameters.Login.Endpoint

		c = &http.Client{Jar: jar}
		_, err = c.PostForm(loginUrl.String(), url.Values{
			configParameters.Login.Fields.getUsernameField(): {
				configParameters.Login.Fields.getUsernameValue()},
			configParameters.Login.Fields.getPasswordField(): {
				configParameters.Login.Fields.getPasswordValue()},
		})
	}

	return err
}

func logout(c *http.Client, environment int) error {
	var err error

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)
	if nil == err {
		logoutUrl := new(url.URL)
		logoutUrl.Scheme = "http"
		logoutUrl.Host = configParameters.Domains[environment]
		logoutUrl.Path = configParameters.Login.Endpoint

		c = &http.Client{Jar: jar}
		_, err = c.PostForm(logoutUrl.String(), url.Values{})
	}

	return err
}

func diffEndpoint(endpoint string) {
	templateUrl := new(url.URL)

	templateUrl.Scheme = "http"
	templateUrl.RawPath = endpoint

	templateUrl.Host = configParameters.Domains[0]

	log.Printf("url: %s", endpoint)
	respPrd, err := clientProduction.Get(templateUrl.String())
	if nil != err {
		log.Println("Error: ", err)
	} else {
		templateUrl.Host = configParameters.Domains[1]

		respStg, err := clientProduction.Get(templateUrl.String())
		if nil != err {
			log.Println("Error: ", err)
		} else {
			prdContent, _ := ioutil.ReadAll(respPrd.Body)
			stgContent, _ := ioutil.ReadAll(respStg.Body)

			log.Printf("prd content: %s\n", prdContent)
			log.Printf("stg content: %s\n", stgContent)
		}
	}
}
