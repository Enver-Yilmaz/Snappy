package main

import "log"


//	Contains Snappy specific menus, really just generic menus like search and settings


//the main menu
func MainMenu(){
	drawLogo.Store(true)
	go ClearAndAppend(
		NewMenuItem("Search", func(){
			drawLogo.Store(false)
			SearchMenu()
		}),
		NewMenuItem("Twitch", TwitchSearchMenu),
		NewMenuItem("Youtube", YoutubeMenu),
		NewMenuItem("Settings", func(){
			log.Println("settings")
		}))
}

//the search menu
func SearchMenu(){
	go ClearAndAppend(
		NewMenuItem("The Movie Database", TMDBMenu),
		NewMenuItem("Alluc", AllucSearchMenu),
		NewMenuItem("Back", MainMenu))
}
