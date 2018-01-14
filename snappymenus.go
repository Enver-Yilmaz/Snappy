package main

import "log"


//	Contains Snappy specific menus, really just generic menus like search and settings


//the main menu
func MainMenu(){
	go ClearAndAppend(
		NewMenuItem("Search", func(){
			SearchMenu()
		}),
		NewMenuItem("Settings", func(){
			log.Println("settings")
		}))
}

//the search menu
func SearchMenu(){
	go ClearAndAppend(
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
