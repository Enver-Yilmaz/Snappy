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
		NewMenuItem("Settings", func(){
			log.Println("settings")
		}))
}

//the search menu
func SearchMenu(){
	go ClearAndAppend(
		NewMenuItem("The Movie Database", func(){
			TMBDMenu()
		}),
		NewMenuItem("Alluc", func(){
			AllucSearchMenu()
		}),
		NewMenuItem("Twitch", func(){
			TwitchSearchMenu()
		}),
		NewMenuItem("Back", func(){
			MainMenu()
		}))
}
