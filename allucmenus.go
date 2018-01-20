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
	go InputMenu("Enter a show/movie that you want to find e.g. the simpsons s04e11", func(){
		SetInputMode(false)	//definitely set this BEFORE we start searching otherwise Snappy just looks weird... Maybe add a loading screen?
		AllucResultMenu(AllucSearchMenu)
	}, SearchMenu)
}

//displays the results of an alluc search
func AllucResultMenu(returnTo func()){
	go ClearMenu()
	//note that the search length is most optimal in multiples of 4 given the rpi3 quad core cpu

	resp, err := http.Get("https://www.alluc.ee/api/search/stream/?apikey=" + allucKey + "&query=" + url.QueryEscape(inputText.Load() + " host:" + strings.Join(hosters, ",")) + "&count=" + strconv.Itoa(AllucSearchLength) + "&from=0&getmeta=0")
	log.Println("https://www.alluc.ee/api/search/stream/?apikey=" + allucKey + "&query=" + url.QueryEscape(inputText.Load() + " host:" + strings.Join(hosters, ",")) + "&count=" + strconv.Itoa(AllucSearchLength) + "&from=0&getmeta=0")

	if err != nil {
		go ClearAndAppend(NewMenuItem("couldn't complete the alluc API request, check your connection and API keys (click to return)", func(){
			returnTo()
		}))
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		termbox.Close()
		log.Fatal(err)
	}

	log.Println(string(body))

	if status, err := jsonparser.GetString(body, "status"); status == "error" || err != nil {
		go ClearAndAppend(
			NewMenuItem("something went wrong :( Maybe there weren't any results? (click to return)", func() {
				returnTo()
			}))
		return
	}

	go AppendMenu(NewMenuItem("Return to search", func(){
		returnTo()
	}))

	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error){

		title, err := jsonparser.GetString(value, "sourcetitle")

		if err != nil {
			title, err = jsonparser.GetString(value, "title")
			if err != nil {
				title = "Somehow there didn't seem to be any title.."
			}
		}

		hostURL, err := jsonparser.GetString(value, "hosterurls", "[0]", "url") //it doesn't seem like the hosterurls array ever has more than one element but i'll keep an eye on the logs

		if err != nil {
			log.Printf("The following alluc result was missing a hosterurl?: %s", string(body))
		} else {
			//TODO refactor this to follow the issue on github
			unrestrict(hostURL, title)
/*			if containsHoster(hostURL, hosters) || strings.Contains(hostURL, "rapidvideo") { //I think this will prevent some time wasting on false links (thanks Alluc) :(
				unrestrict(hostURL, title)
			} else {
				log.Printf("The host URL(%s) didn't contain any known hoster", hostURL)
			}*/
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