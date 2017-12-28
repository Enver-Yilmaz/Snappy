package main

import (
	"github.com/nsf/termbox-go"
	"log"
	"os"
	"net/http"
	"strconv"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/url"
	"os/exec"
	"path/filepath"
)

/*
	Snappy!
*/

//TODO this note might be written somewhere else already, but I'd like to implement sqlite caching

var(
	menu = NewMenu([]MenuItem{})
	inputText = ""
	inPlayback = false //todo: do something with this maybe
	desktop = false
	drawChan = make(chan int)
	allucKey = ""
	realDebridKey = ""
	tmdbKey = ""
	exPath = ""
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

	mainMenu()
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
func mainMenu(){
	go ClearAndAppend(
		NewMenuItem("Search", func(){
		searchMenu()
		}),
		NewMenuItem("Settings", func(){
			log.Println("settings")
		}))
}

//the search menu
func searchMenu(){
	go ClearAndAppend(
		NewMenuItem("Alluc", func(){
			allucSearchMenu()
		}),
		NewMenuItem("Back", func(){
			mainMenu()
		}))
}

//the t.v. search menu
//really just the alluc search menu, some inappropriate naming going on here
//TODO i just started changing this but didn't finish the names that is
func allucSearchMenu(){
	go ClearAndAppend(
		NewMenuItem("Enter", func(){
			if inputText != ""{
				allucResultMenu()
			}
		}),
		NewMenuItem("Back", func(){
			searchMenu()
			go SetInputMode(false)
		}))
	go SetInputMode(true)
}

//displays the results of an alluc search
func allucResultMenu(){
	var searchLength = 8 //TODO read this from a config file or something
	go SetInputMode(false)
	go ClearMenu()
	//note that the search length is most optimal in multiples of 4 given the rpi3 quad core cpu

	resp, err := http.Get("https://www.alluc.ee/api/search/stream/?apikey=" + allucKey + "&query=" + url.QueryEscape(inputText + " host:openload.co,thevideo.me,bitporno.sx,cloudtime.to,dailymotion.com,datoporn.com,flashx.tv,wholecloud.net,novamov.com,auroravid.to,nowvideo.co,rapidvideo.ws,redtube.com,userscloud.com,videoweed.com,vimeo.com,youporn.com") + "&count=" + strconv.Itoa(searchLength) +"&from=0&getmeta=0")
	check(err)

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	check(err)

	log.Println(string(body))

	if status, err := jsonparser.GetString(body, "status"); status == "error" || err != nil {
		go ClearAndAppend(
			NewMenuItem("something went wrong :( Maybe there weren't any results? (click to return)", func() {
				allucSearchMenu()
			}))
		return
	}

	//TODO split this into seperate threads rather than for loop style
	//split into chunks of 4 if possible and thread out from there
	//given I think this is pinning the cpu most most right now, because we're blocking on real-debrid unrestricting

	go AppendMenu(NewMenuItem("Return to search", func(){
		allucSearchMenu()
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
	//blockChan <- 1
	go func(){
		//defer func(){
			//<-blockChan
		//}()

		form := url.Values{}
		form.Add("link", link)

		resp, err := http.PostForm("https://api.real-debrid.com/rest/1.0/unrestrict/link?auth_token="+realDebridKey, form)

		check(err)
		body, err := ioutil.ReadAll(resp.Body)

		check(err)

		log.Println(string(body))

		streamURL, err := jsonparser.GetString(body, "download")

		if err != nil {
			log.Println("couldn't debrid unrestrict this link: " + link)
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
				inPlayback = true
			}))
		}

	}()
}



/*
	This is a very simple set up to allow control of Snappy from a HTTP request, it simply parses the queries and calls the respective menu functions
	It's very likely to cause unexpected behavior if a user is inputting with something like a keyboard and the remote at the same time
	Because the remote is asynchronously changing the menu, but it should be fine for now because that use case is unlikely (said every programmer ever)
	//TODO cache things like originally planned and add a back button

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
			inputText = text
		}
	} else if v.Get("home") != ""{
		isReturn, err := strconv.ParseBool(v.Get("home"))
		check(err)
		if isReturn{
			mainMenu()
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
