package main

import (
	"flag"
	"fmt"
	"net/http"
)

//We pack the data up into a struct
//This is one way to make sure that the data arrives together at the same time
type webData struct {
	rawArg string
	status string
}

func (results *webData) String() string {
	return "Status " + results.status + "for " + results.rawArg
}

// Major function that gets and prepares the data based on flags
func grabWeb(url string, status, quiet bool, data chan webData, quit chan bool) {

	//Notice that magic is the results of the Get request
	//err stores a returned error code from the Get request if it occurs
	//err is nil otherwise
	magic, err := http.Get(url)

	//error handling
	if err != nil {
		if !quiet {
			fmt.Println("Failed to proccess:", url)
		}
		quit <- true
		return
	}
	results := webData{rawArg: url}
	if status {
		results.status = magic.Status
	}
	magic.Body.Close()

	//Order of channel use matters here
	data <- results
	quit <- true

}

func main() {
	//Command line flags
	nothingFlag := flag.Bool("nothing", false, "Program does nothing, default is false")
	statusFlag := flag.Bool("status", true, "Get server status, default is true")
	quietFlag := flag.Bool("quiet", false, "Don't print failures, default is false")

	//Flags need to be parsed in order to be usable
	flag.Parse()

	//Quit early if no arguements are present or nothingFlag is set
	if len(flag.Args()) == 0 || *nothingFlag {
		return
	}

	//Our channels
	quit := make(chan bool)
	data := make(chan webData)

	//Iterate over the cmd line arguements
	for _, value := range flag.Args() {
		go grabWeb(value, *statusFlag, *quietFlag, data, quit)
	}

	//Process responses from the goroutines
	for counter := 0; counter < len(flag.Args()); {
		select {
		case results := <-data:
			if *statusFlag {
				fmt.Println(results.String())
			}
		case <-quit:
			counter++
		}
	}
}
