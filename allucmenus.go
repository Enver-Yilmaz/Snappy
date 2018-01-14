package main

import (
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"io/ioutil"
	"log"
	"github.com/buger/jsonparser"
	"github.com/nsf/termbox-go"
)

//this file contains menus specific to Alluc

//the t.v. search menu
func AllucSearchMenu(){
	go ClearAndAppend(
		NewMenuItem("Enter", func(){
			if inputText.Load() != ""{
				AllucResultMenu()
				go SetInputMode(false)
			}
		}),
		NewMenuItem("Back", func(){
			go SetInputMode(false)
			SearchMenu()
		}))
	go SetInputMode(true)
}

//displays the results of an alluc search
func AllucResultMenu(){
	go SetInputMode(false)
	go ClearMenu()
	//note that the search length is most optimal in multiples of 4 given the rpi3 quad core cpu

	resp, err := http.Get("https://www.alluc.ee/api/search/stream/?apikey=" + allucKey + "&query=" + url.QueryEscape(inputText.Load()) + " host:" + strings.Join(hosters, ",") + "&count=" + strconv.Itoa(AllucSearchLength) + "&from=0&getmeta=0")

	if err != nil {
		go ClearAndAppend(NewMenuItem("couldn't complete the alluc API request, check your connection and API keys (click to return)", func(){
			AllucSearchMenu()
		}))
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		termbox.Close()
		log.Fatal(err)
	}

	if status, err := jsonparser.GetString(body, "status"); status == "error" || err != nil {
		go ClearAndAppend(
			NewMenuItem("something went wrong :( Maybe there weren't any results? (click to return)", func() {
				AllucSearchMenu()
			}))
		return
	}

	go AppendMenu(NewMenuItem("Return to search", func(){
		AllucSearchMenu()
	}))

	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error){

		title, err := jsonparser.GetString(value, "sourcetitle")

		if err != nil {
			title, err = jsonparser.GetString(value, "title")
			if err != nil {
				title = "Somehow there didn't seem to be any title.."
			}
		}

		hostURL, err := jsonparser.GetString(value, "hosterurls", "[0]", "url")

		if err != nil {
			log.Printf("The following alluc result was missing a hosterurl?: %s", string(body))
		} else {
			if containsHoster(hostURL, hosters){ //I think this will prevent some time wasting on false links (thanks Alluc) :(
				unrestrict(hostURL, title)
			} else {
				log.Printf("The host URL(%s) didn't contain any known hoster", hostURL)
			}
		}
	}, "result")
}


//checks if the given url contains any of the hosters
func containsHoster(url string, hosters []string)bool{
	for _, hoster := range hosters {
		if strings.Contains(url, hoster){
			return true
		}
	}
	return false
}