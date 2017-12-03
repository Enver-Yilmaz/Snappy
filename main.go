package main

import (
	"github.com/nsf/termbox-go"
	"github.com/mattn/go-runewidth"
	"log"
	"os"
	"net/http"
	"strconv"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/url"
	"os/exec"
)

/*
	Snappy!
*/

const allucKey = "???"
const realDebridKey = "???"
//TODO load these keys from a config too


var menu  = Menu{nil, 0, false}
var inputText = ""

func main() {
	err := termbox.Init()
	check(err)

	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)

	err = os.Remove("log")
	check(err)

	f, err := os.OpenFile("log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	log.SetOutput(f)

	http.HandleFunc("/", remoteHandler)
	go http.ListenAndServe(":8080", nil)

	mainMenu()
	drawAll()

	mainloop:
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key{
				case termbox.KeyArrowDown:
					menu.CursorDown()
				case termbox.KeyArrowUp:
					menu.CursorUp()
				case termbox.KeyCtrlC:
					break mainloop
				case termbox.KeyEnter:
					menu.MenuSelect()
				}
			case termbox.EventResize:
				drawAll()
			case termbox.EventError:
				panic(ev.Err)
			}
			drawAll()
		}

}

func drawAll(){
	w, _ := termbox.Size()

	logo := []string{
		"     _______..__   __.      ___      .______   .______   ____    ____ ",
		"    /       ||  \\ |  |     /   \\     |   _  \\  |   _  \\  \\   \\  /   / ",
		"   |   (----`|   \\|  |    /  ^  \\    |  |_)  | |  |_)  |  \\   \\/   /  ",
		"    \\   \\    |  . `  |   /  /_\\  \\   |   ___/  |   ___/    \\_    _/   ",
		".----)   |   |  |\\   |  /  _____  \\  |  |      |  |          |  |     ",
		"|_______/    |__| \\__| /__/     \\__\\ | _|      | _|          |__|     ",
		"                                                                      ",
	}

	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	//draw the logo
	//TODO: add condition to hide
	for i, line := range logo {
		tbprint(w / 2 - len(line) / 2, i, termbox.ColorGreen, termbox.ColorBlack, line)
	}

	menu.DrawMenu(len(logo) + 1)

	termbox.Flush()
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

//the main menu
func mainMenu(){
	menu = NewMenu([]MenuItem{
		NewMenuItem("Search", func(){
			searchMenu()
		}),
		NewMenuItem("Settings", func(){
			log.Println("settings")
		}),
	})
}

//the search menu
func searchMenu(){
	menu = NewMenu([]MenuItem{
		NewMenuItem("T.V. show", func(){
			tvSearchMenu()
		}),
		NewMenuItem("Back", func(){
			mainMenu()
		}),
	})
}

//the t.v. search menu
func tvSearchMenu(){
	menu = NewMenu([]MenuItem{
		NewMenuItem("Enter", func(){
			if inputText != ""{
				allucSearchMenu()
			}
		}),
		NewMenuItem("Back", func(){
			searchMenu()
		}),
	})
	menu.inputMode = true
}

/*
*	IT WORKS!!!
*	It's a little janky but it works!
*/
func allucSearchMenu(){
	var searchLength = 2 //TODO read this from a config file or something

	resp, err := http.Get("https://www.alluc.ee/api/search/stream/?apikey=" + allucKey + "&query=" + url.QueryEscape(inputText + " host:openload.co,thevideo.me,") + "&count=" + strconv.Itoa(searchLength) +"&from=0&getmeta=0")
	check(err)

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	check(err)

	log.Println(string(body))

	var menuItems = make([]MenuItem, searchLength)

	i := 0
	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error){

		title, err := jsonparser.GetString(value, "sourcetitle")
		log.Println(title)
		log.Println()

		check(err)

		hostURL, err := jsonparser.GetString(value, "hosterurls", "[0]", "url")
		log.Println(hostURL)
		check(err)

		form := url.Values{}
		form.Add("link", hostURL)

		resp, err := http.PostForm("https://api.real-debrid.com/rest/1.0/unrestrict/link?auth_token="+realDebridKey, form)

		check(err)

		body, err := ioutil.ReadAll(resp.Body)

		check(err)

		streamURL, err := jsonparser.GetString(body, "download")

		log.Println(streamURL)

		menuItems[i] = NewMenuItem(title, func() {

			out, err := exec.Command("vlc", streamURL).Output()

			if err != nil {
				panic(err)
			}

			log.Println(out)

		})

		i++
	}, "result")

	menu = NewMenu(menuItems)
}

/*
	This is a very simple set up to allow control of Snappy from a HTTP request, it simply parses the queries and calls the respective menu functions
	It's very likely to cause unexpected behavior if a user is inputting with something like a keyboard and the remote at the same time
	Because the remote is asynchronously changing the menu, but it should be fine for now because that use case is unlikely (said every programmer ever)
	//TODO likely find a smarter way to handle these race conditions
	//TODO cache things like originally planned and add a back button

	current commands:
	down(int) move the menu down the amount e.g. ?down=1
	up(int) move the menu up the amount e.g. ?up=1
	select(bool) selects the current button if true e.g. ?select=true I'm not sure why this is even a bool tbh
	text(string) sets the input text to the given string if the menu is in input mode e.g. text?=spongebob
*/
func remoteHandler(w http.ResponseWriter, r *http.Request){
	v := r.URL.Query()

	if v.Get("down") != ""{
		down, err := strconv.Atoi(v.Get("down"))
		check(err)
		for i := 0; i < down; i++{
			menu.CursorDown()
		}
	} else if v.Get("up") != ""{
		up, err := strconv.Atoi(v.Get("up"))
		check(err)
		for i := 0; i < up; i++{
			menu.CursorUp()
		}
	} else if v.Get("select") != ""{
		isSelected, err := strconv.ParseBool(v.Get("select"))
		check(err)
		if isSelected{
			menu.MenuSelect()
		}
	} else if v.Get("text") != ""{
		text := v.Get("text")
		if menu.inputMode{
			inputText = text
		}
	} else if v.Get("home") != ""{
		isReturn, err := strconv.ParseBool(v.Get("home"))
		check(err)
		if isReturn{
			mainMenu()
		}
	}

	drawAll()
	w.WriteHeader(http.StatusOK)
}


func check(err error){
	if err != nil {
		panic(err)
	}
}
