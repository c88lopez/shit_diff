package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"io/ioutil"

	"encoding/json"

	"encoding/csv"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/publicsuffix"
)

var clientProduction http.Client
var clientStage http.Client

var configParameters configJson

var results [2]string
var rows [][]string

func main() {
	// Login to ruby
	log.Println("Start...")

	bootstrap()

	log.Println("Login to production...")
	if err := login(0); nil != err {
		log.Fatal(err)
	}
	defer logout(0)

	log.Println("Login to stage...")
	if err := login(1); nil != err {
		log.Fatal(err)
	}
	defer logout(1)

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

func login(environment int) error {
	var err error
	//var resp *http.Response

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)
	if nil != err {
		return err
	}

	loginUrl := new(url.URL)
	loginUrl.Scheme = "http"
	loginUrl.Host = configParameters.Domains[environment]
	loginUrl.Path = configParameters.Login.Endpoint

	if 0 == environment {
		clientProduction = http.Client{Jar: jar}
		_, err = clientProduction.PostForm(loginUrl.String(), url.Values{
			configParameters.Login.Fields.getUsernameField(): {
				configParameters.Login.Fields.getUsernameValue()},
			configParameters.Login.Fields.getPasswordField(): {
				configParameters.Login.Fields.getPasswordValue()},
		})

		if nil != err {
			log.Fatal(err)
		}
	} else {
		clientStage = http.Client{Jar: jar}
		_, err = clientStage.PostForm(loginUrl.String(), url.Values{
			configParameters.Login.Fields.getUsernameField(): {
				configParameters.Login.Fields.getUsernameValue()},
			configParameters.Login.Fields.getPasswordField(): {
				configParameters.Login.Fields.getPasswordValue()},
		})

		if nil != err {
			log.Fatal(err)
		}
	}

	return err
}

func logout(environment int) error {
	var err error

	if nil == err {
		logoutUrl := new(url.URL)
		logoutUrl.Scheme = "http"
		logoutUrl.Host = configParameters.Domains[environment]
		logoutUrl.Path = configParameters.Login.Endpoint

		if 0 == environment {
			_, err = clientProduction.PostForm(logoutUrl.String(), url.Values{})
		} else {
			_, err = clientStage.PostForm(logoutUrl.String(), url.Values{})
		}
	}

	return err
}

func diffEndpoint(endpoint string) {
	urlTool := new(url.URL)

	fullUrl, err := urlTool.Parse("http://" + configParameters.Domains[0] + "/" + endpoint)
	if nil != err {
		log.Fatal(err)
	}

	respPrd, err := clientProduction.Get(fullUrl.String())
	if nil != err {
		log.Println("Error: ", err)
	} else {
		fullUrl.Host = configParameters.Domains[1]

		respStg, err := clientStage.Get(fullUrl.String())
		if nil != err {
			log.Println("Error: ", err)
		} else {
			prdContent, err := ioutil.ReadAll(respPrd.Body)
			if nil != err {
				log.Fatal(err)
			}

			stgContent, err := ioutil.ReadAll(respStg.Body)
			if nil != err {
				log.Fatal(err)
			}

			log.Printf("prd content: %s\n", prdContent)
			log.Printf("stg content: %s\n", stgContent)
		}
	}
}
