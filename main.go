package main

import (
	"github.com/nsf/termbox-go"
	"go.uber.org/atomic"
	"log"
	"os"
	"net/http"
	"path/filepath"
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
	inputMenuTitle = atomic.NewString("")
	inPlayback = atomic.NewBool(false)
	drawLogo = atomic.NewBool(true)
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

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func main() {
	defer termbox.Close() //I'm pretty sure this still has to be deferred in main unfortunately

	if Exists(exPath + "/log"){
		err := os.Remove(exPath + "/log")
		if err != nil {
			log.Println(err.Error())	//I doubt this error will every actually happen, but if it does we probably shouldn't panic right?
		}
	}

	f, err := os.OpenFile(exPath + "/log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err == nil {
		defer f.Close()
		log.SetOutput(f)
	}

	http.HandleFunc("/", RemoteHandler)
	go http.ListenAndServe(":8080", nil)

	MainMenu()

	go menu.TBdraw()
	drawChan <- 1
	TBinput()
}
