package main

import (
	"github.com/nsf/termbox-go"
	"go.uber.org/atomic"
	"log"
	"os"
	"net/http"
	"strconv"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
	Snappy!
*/

//TODO this note might be written somewhere else already, but I'd like to implement sqlite caching
//probably still do

const desktop = true

var(
	menu = NewMenu([]MenuItem{})
	inputText = atomic.NewString("")
	inPlayback = atomic.NewBool(false)
	drawChan = make(chan int)
	//the keys aren't atomic because they shouldn't ever change and concurrent reads should be fine
	allucKey = ""
	realDebridKey = ""
	tmdbKey = ""
	exPath = ""
	AllucSearchLength = 4 //default ?
	hosters = []string {
		"openload.co",
		"thevideo.me",
		"bitporno.sx",
		"cloudtime.to",
		"datoporn.com",
		"flashx.tv",
		"wholecloud.net",
		"novamov.com",
		"auroravid.to",
		"rapidvideo.ws",
		"redtube.com",
		"userscloud.com",
		"youporn.com",
	}
)

//init here just serves the purpose of unpacking the config file initializing the variables associated with it
//also inits the log
func init(){
	ex, err := os.Executable()

	if err != nil {
		panic(err)
	}

	exPath = filepath.Dir(ex)

	config := ParseConfigFile(exPath)
	allucKey = config.AllucKey
	realDebridKey = config.RealDebridKey
	tmdbKey = config.TmdbKey
	AllucSearchLength = config.AllucSearchLength
}

func main() {
	defer termbox.Close() //I'm pretty sure this still has to be deferred in main unfortunately

	if Exists(exPath + "/log"){
		err := os.Remove(exPath + "/log")
		check(err)
	}

	f, err := os.OpenFile(exPath + "/log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err == nil {
		defer f.Close()
		log.SetOutput(f)
	}

	http.HandleFunc("/", remoteHandler)
	go http.ListenAndServe(":8080", nil)

	MainMenu()

	go menu.TBdraw()
	drawChan <- 1
	TBinput()
}


// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

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

func TwitchSearchMenu(){
	go ClearAndAppend(
		NewMenuItem("Enter", func(){
			if inputText.Load() != "" {
				TwitchStreamsMenu(inputText.Load())
				go SetInputMode(false)
			}
		}),
		NewMenuItem("Back", func(){
			MainMenu()
			go SetInputMode(false)
		}),
	)
	go SetInputMode(true)
}

//displays the results of an alluc search
func AllucResultMenu(){
	go SetInputMode(false)
	go ClearMenu()
	//note that the search length is most optimal in multiples of 4 given the rpi3 quad core cpu

	resp, err := http.Get("https://www.alluc.ee/api/search/stream/?apikey=" + allucKey + "&query=" + url.QueryEscape(inputText.Load()) + " host:" + strings.Join(hosters, ",") + "&count=" + strconv.Itoa(AllucSearchLength) + "&from=0&getmeta=0")
	check(err)

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	check(err)

	log.Println(string(body))

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
			check(err)
		}
		log.Println(title)

		hostURL, err := jsonparser.GetString(value, "hosterurls", "[0]", "url")
		log.Println(hostURL)
		check(err)

		unrestrict(hostURL, title)


	}, "result")
}

func unrestrict(link string, title string){
	go func(){

		form := url.Values{}
		form.Add("link", link)

		resp, err := http.PostForm("https://api.real-debrid.com/rest/1.0/unrestrict/link?auth_token="+realDebridKey, form)

		check(err)
		body, err := ioutil.ReadAll(resp.Body)

		check(err)

		streamURL, err := jsonparser.GetString(body, "download")

		if err != nil {
			log.Println("couldn't debrid unrestrict this link: " + link)
			log.Println("Error body:" + string(body))
		} else {
			log.Println(streamURL)

			go AppendMenu(NewMenuItem(title, func() {

				if !desktop{
					command := exec.Command("omxplayer", "-b", "-o", "hdmi", "--live", streamURL)
					err = command.Run()
				} else {
					command := exec.Command("vlc", streamURL)
					err = command.Run()
				}
				inPlayback.Store(true)
			}))
		}

	}()
}



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
func remoteHandler(w http.ResponseWriter, r *http.Request){
	v := r.URL.Query()

	if v.Get("down") != ""{
		down, err := strconv.Atoi(v.Get("down"))
		check(err)
		for i := 0; i < down; i++{
			go CursorDown()
		}
	} else if v.Get("up") != ""{
		up, err := strconv.Atoi(v.Get("up"))
		check(err)
		for i := 0; i < up; i++{
			go CursorUp()
		}
	} else if v.Get("select") != ""{
		isSelected, err := strconv.ParseBool(v.Get("select"))
		check(err)
		if isSelected{
			go MenuSelect()
		}
	} else if v.Get("text") != ""{
		text := v.Get("text")
		if menu.inputMode{
			inputText.Store(text)
		}
	} else if v.Get("home") != ""{
		isReturn, err := strconv.ParseBool(v.Get("home"))
		check(err)
		if isReturn{
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


func check(err error){
	if err != nil {
		panic(err)
	}
}
