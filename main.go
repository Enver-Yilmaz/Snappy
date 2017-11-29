package main

import (
	"github.com/nsf/termbox-go"
	"github.com/mattn/go-runewidth"
	"log"
	"os"
)

var menu  = Menu{nil, 0}

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

func searchMenu(){
	menu = NewMenu([]MenuItem{
		NewMenuItem("T.V. show", func(){
			log.Println("search for t.v. show")
		}),
		NewMenuItem("main menu", func(){
			mainMenu()
		}),
	})
}


func check(err error){
	if err != nil {
		panic(err)
	}
}