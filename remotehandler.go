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
	down(int) move the menu down the amount e.g. ?down=1
	up(int) move the menu up the amount e.g. ?up=1
	select(bool) selects the current button if true e.g. ?select=true I'm not sure why this is even a bool tbh
	text(string) sets the input text to the given string if the menu is in input mode e.g. ?text=spongebob
	home(bool) returns home if true e.g. ?home=true
	omxcommand(someomxcommand) executes the command for omxplayer via dbus through a bash script... e.g. ?omxcommand=stop
*/

//the actual function for remote handling
func RemoteHandler(w http.ResponseWriter, r *http.Request){
	v := r.URL.Query()

	if v.Get("down") != ""{
		down, err := strconv.Atoi(v.Get("down"))

		if err != nil {
			log.Println(err.Error())
			return	//something went wrong? we'll just log it and ignore it, I don't think it really matters on something like the remote handler
		}

		for i := 0; i < down; i++{
			go CursorDown()
		}
	} else if v.Get("up") != ""{
		up, err := strconv.Atoi(v.Get("up"))

		if err != nil {
			log.Println(err.Error())
			return
		}

		for i := 0; i < up; i++{
			go CursorUp()
		}
	} else if v.Get("select") != ""{
		isSelected, err := strconv.ParseBool(v.Get("select"))

		if err != nil {
			log.Println(err.Error())
			return
		}

		if isSelected{
			go MenuSelect()
		}
	} else if v.Get("text") != ""{
		text := v.Get("text")
		if menu.inputMode{
			inputText.Store(text) //this could be dangerous?
		}
	} else if v.Get("home") != ""{
		isReturn, err := strconv.ParseBool(v.Get("home"))

		if err != nil {
			log.Println(err.Error())
			return
		}

		if isReturn { //yeah this is really dumb
			MainMenu()
			go SetInputMode(false)
		}
	} else if v.Get("omxcommand") != ""{
		command := v.Get("omxcommand")
		log.Println("command was " + command)
		exec.Command("/root/omxdbus.sh", command).Run() //this probably isn't the smartest way to do this, but I don't think it really matters
	}

	drawChan <- 1

	w.WriteHeader(http.StatusOK)
}