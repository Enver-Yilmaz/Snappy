package main

import (
	"net/http"
	"strconv"
	"log"
	"os/exec"
)

//TODO consider using routes instead of weird query params

/*
	This is a very simple set up to allow control of Snappy from a HTTP request, it simply parses the queries and calls the respective menu functions
	It's very likely to cause unexpected behavior if a user is inputting with something like a keyboard and the remote at the same time
	Because the remote is asynchronously changing the menu, but it should be fine for now because that use case is unlikely (said every programmer ever)

	current commands:
	snappycommand(somesnappycommand) executes the given snappy command e.g. ?snappycommand=home
	text(sometext) sets the inputtext variable to the given text e.g. ?text=vinesauce
	omxcommand(someomxcommand) executes the command for omxplayer via dbus through a bash script... e.g. ?omxcommand=stop
*/

//the actual function for remote handling
func RemoteHandler(w http.ResponseWriter, r *http.Request){
	v := r.URL.Query()

	if command := v.Get("snappycommand"); command != "" {
		switch command {
		case "down":
			go CursorDown()
		case "up":
			go CursorUp()
		case "select":
			go MenuSelect()
		case "home":
			go MainMenu()
		}
	} else if v.Get("text") != ""{
		text := v.Get("text")
		if menu.inputMode{
			inputText.Store(text)
		}
	} else if v.Get("omxcommand") != ""{
		command := v.Get("omxcommand")
		log.Println("command was " + command)
		exec.Command("/root/omxdbus.sh", command).Start() //this probably isn't the smartest way to do this, but I don't think it really matters
	}

	drawChan <- 1

	w.WriteHeader(http.StatusOK)
}