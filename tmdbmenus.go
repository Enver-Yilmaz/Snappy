package main

import (
	"net/http"
	"fmt"
	"net/url"
	"log"
	"io/ioutil"
	"github.com/nsf/termbox-go"
	"github.com/buger/jsonparser"
)

const tmdbapi = "https://api.themoviedb.org/3/"


//Generic TMDB menu that splits out into tv and movie search
func TMBDMenu(){
	go ClearAndAppend(
		NewMenuItem("T.V. Show", TMDBTVSearch),
		NewMenuItem("Movie", TMDBMovieSearch),
		NewMenuItem("Back", SearchMenu))
}

//allows the user to search TMDB for a t.v. show
func TMDBTVSearch(){
	go InputMenu("Enter the name of a T.V. show that you want to find e.g. Seinfeld", func(){
		TMDBSearch("search/tv", inputText.Load())
	}, TMBDMenu)
}

//allows the user to search TMDB for a movie
func TMDBMovieSearch(){
	go InputMenu("Enter the name of a movie that you want to find e.g. Requiem for a Dream", func(){
		TMDBSearch("search/movie", inputText.Load())
	}, TMBDMenu)
}

//TODO maybe add pages
//this function searches tmdb, and builds a menu with the results, with the according route (tv/movie/more?)
func TMDBSearch(route string, searchString string){
	var finalPart = fmt.Sprintf("?api_key=%s&language=en-US&query=%s", tmdbKey, url.QueryEscape(searchString))

	resp, err := http.Get(tmdbapi + route + finalPart)
	log.Println(tmdbapi + route + finalPart)

	if err != nil {
		ClearAndAppend(NewMenuItem("Failed to connect to the TMDB API (click to return)", TMBDMenu))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		termbox.Close()
		log.Fatal(err)
	}

	log.Println(string(body))

	ClearMenu()

	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error){
		if route == "search/tv" {

			name, err := jsonparser.GetString(value, "original_name")

			if err != nil {
				termbox.Close()
				log.Fatal(err)
			}

			id, err := jsonparser.GetInt(value, "id")

			if err != nil {
				termbox.Close()
				log.Fatal(err)
			}

			go AppendMenu(NewMenuItem(name, func(){
				TMDBTVSeasonMenu(int(id), name)
			}))

		} else {
			name, err := jsonparser.GetString(value, "title")

			if err != nil {
				termbox.Close()
				log.Fatal(err) //TODO find a better way to handle this error
			}

			go AppendMenu(NewMenuItem(name, func(){
				inputText.Store(name)
				AllucResultMenu(TMDBMovieSearch)
			}))
		}
	}, "results")

	if route == "search/tv" {
		AppendMenu(NewMenuItem("return to search", TMDBTVSearch))
	} else {
		AppendMenu(NewMenuItem("return to search", TMDBMovieSearch))
	}
}

//builds a menu for the given TMDB ID
func TMDBTVSeasonMenu(tvId int, name string) {
	finalPart := fmt.Sprintf("tv/%d?api_key=%s&language=en-US", tvId, tmdbKey)
	requestURI := tmdbapi + finalPart

	log.Print(requestURI)

	resp, err := http.Get(requestURI)

	if err != nil {
		ClearAndAppend(NewMenuItem("Failed to connect to the TMDB API (click to return)", TMDBTVSearch))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		termbox.Close()
		log.Fatal(err)
	}

	log.Println(string(body))

	ClearMenu()

	AppendMenu(NewMenuItem("Return to TMDB TV search", func(){
		ClearMenu() //TODO fix all menus to clear at the beginning, this is just dumb
		TMDBSearch("search/tv", inputText.Load()) //hopefully the input text didn't get changed...
	}))

	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error){

		seasonNumber, err := jsonparser.GetInt(value, "season_number")

		if err != nil {
			termbox.Close()
			log.Fatal(err)
		}

		episodeCount, err := jsonparser.GetInt(value, "episode_count")

		if err != nil {
			termbox.Close()
			log.Fatal(err)
		}

		var title = ""

		if int(seasonNumber) == 0 {
			title = fmt.Sprintf("%d Special Episodes", int(episodeCount))
		} else {
			title = fmt.Sprintf("Season %d with %d episodes", int(seasonNumber), int(episodeCount))
		}

		AppendMenu(NewMenuItem(title, func(){
			TMDBTVEpisodeMenu(tvId, int(seasonNumber), name)
		}))

	}, "seasons")
}

func TMDBTVEpisodeMenu(tvId, seasonNumber int, showname string) {
	finalPart := fmt.Sprintf("tv/%d/season/%d?api_key=%s&lanuage=en-US", tvId, seasonNumber, tmdbKey)
	requestURI := tmdbapi + finalPart

	log.Println(requestURI)

	resp, err := http.Get(requestURI)

	if err != nil {
		ClearAndAppend(NewMenuItem("Failed to connect to the TMDB API (click to return)", TMDBTVSearch))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		termbox.Close()
		log.Fatal(err)
	}

	log.Println(string(body))

	ClearMenu()

	AppendMenu(NewMenuItem("Return to Seasons", func(){
		ClearMenu()
		TMDBTVSeasonMenu(tvId, showname)
	}))

	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error){

		episodeNumber, err := jsonparser.GetInt(value, "episode_number")

		if err != nil {
			termbox.Close()
			log.Fatal(err)
		}

		name, err := jsonparser.GetString(value, "name")

		if err != nil {
			termbox.Close()
			log.Fatal(err)
		}

		title := fmt.Sprintf("Episode %d: %s", int(episodeNumber), name)

		AppendMenu(NewMenuItem(title, func(){
			allucString := fmt.Sprintf("%s season %d episode %d", showname, seasonNumber, int(episodeNumber))
			inputText.Store(allucString)
			AllucResultMenu(func(){
				ClearMenu()
				TMDBTVEpisodeMenu(tvId, seasonNumber, showname)
			})
		}))

	}, "episodes")
}