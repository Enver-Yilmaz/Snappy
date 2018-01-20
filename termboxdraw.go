package main

import (
	"github.com/nsf/termbox-go"
	"strconv"
	"github.com/mattn/go-runewidth"
	"os"
)

var logo = []string{
"     _______..__   __.      ___      .______   .______   ____    ____ ",
"    /       ||  \\ |  |     /   \\     |   _  \\  |   _  \\  \\   \\  /   / ",
"   |   (----`|   \\|  |    /  ^  \\    |  |_)  | |  |_)  |  \\   \\/   /  ",
"    \\   \\    |  . `  |   /  /_\\  \\   |   ___/  |   ___/    \\_    _/   ",
".----)   |   |  |\\   |  /  _____  \\  |  |      |  |          |  |     ",
"|_______/    |__| \\__| /__/     \\__\\ | _|      | _|          |__|     ",
"                                                                      ",
}

func init(){
	if err := termbox.Init(); err != nil {
		panic(err)
	}

	termbox.SetInputMode(termbox.InputEsc)
}

//allows keyboard input to affect snappy while using termbox
func TBinput(){
	inputloop:
		for {
			//todo add mutex for input and probably make the input a struct, so the mutex is part of it
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key{
				case termbox.KeyArrowDown:
					go CursorDown()
				case termbox.KeyArrowUp:
					go CursorUp()
				case termbox.KeyCtrlC:
					break inputloop
				case termbox.KeyEnter:
					MenuSelect()
				case termbox.KeyBackspace:
					if len(inputText.Load()) > 0 {
						//add a backspace here...
					}
				case termbox.KeySpace:
					inputText.Store(inputText.Load() + " ")
				default:
					if menu.inputMode{
						str, err := strconv.Unquote(strconv.QuoteRune(ev.Ch))
						if err != nil {
							panic(err)
						}
						inputText.Store(inputText.Load() + str)
					}
				}
			case termbox.EventResize:
				drawChan <- 1
			case termbox.EventError:
				panic(ev.Err)
			}
			drawChan <- 1
		}
}

//prints a string using termbox-go
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func TBcheck(err error){
	if err != nil {
		termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

		w, h := termbox.Size()

		tbprint(w / 2 - len(err.Error()) / 2, h / 2, termbox.ColorGreen, termbox.ColorBlack, err.Error())
		termbox.Flush()
		termbox.Close()
		os.Exit(-1)
	}
}


//draws snappy via termbox
func (menu *Menu)TBdraw(){
	var yOffset = 0
	for range drawChan {
		menu.lock.Lock()

		//not sure how this happens actually, but oh well
		if menu.currentlySelected > len(menu.Items){
			menu.lock.Unlock()
			continue
		}

		if len(menu.Items) == 0 {
			menu.lock.Unlock()
			continue
		}

		w, h := termbox.Size()

		//todo customization via config
		termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

		if drawLogo.Load() {
			for i, line := range logo {
				tbprint(w / 2 - len(line) / 2, i, termbox.ColorGreen, termbox.ColorBlack, line)
			}
		}

		if !menu.inputMode{
			chunkRadius := (h - yOffset) / 4


			//draw the middle (selected item)
			tbprint(w / 2 - len(menu.Items[menu.currentlySelected].text) / 2, h / 2 + yOffset / 2, termbox.ColorBlue, termbox.ColorBlack, menu.Items[menu.currentlySelected].text)

			//get the lower chunk of other items
			lowerChunk := menu.Items[max(0, menu.currentlySelected - chunkRadius) : menu.currentlySelected]

			i := 0

			//draw them
			for item := range reverse(lowerChunk){
				tbprint(w / 2 - len(item.text) / 2, h / 2 + yOffset / 2 - (i + 1) * 2, termbox.ColorWhite, termbox.ColorBlack, item.text)
				i++
			}

			//get the upper chunk of other items
			upperChunk := menu.Items[menu.currentlySelected + 1 : min(menu.currentlySelected + chunkRadius, len(menu.Items))]

			//draw them
			for i := len(upperChunk) - 1; i >= 0; i-- {
				item := upperChunk[i]
				tbprint(w / 2 - len(item.text) / 2, h / 2 + yOffset / 2 + (i + 1) * 2, termbox.ColorWhite, termbox.ColorBlack, item.text)
			}
		} else { // simple text input, needs to be better in the future TODO future me make this thing better
			text := inputMenuTitle.Load()
			tbprint(w / 2 - len(text) / 2, h / 2 - 5, termbox.ColorWhite, termbox.ColorBlack, text)
			tbprint(w / 2 - len(inputText.Load()) / 2, h / 2, termbox.ColorWhite, termbox.ColorBlack, inputText.Load())
			tbprint(w / 2 - len("Input: ") - len(inputText.Load()) / 2, h / 2, termbox.ColorWhite, termbox.ColorBlack, "Input: ")

			for i, button := range menu.Items{

				if i == menu.currentlySelected{
					tbprint( w / 2 - len(button.text) / 2, h / 2 + 3 + i * 2, termbox.ColorBlue, termbox.ColorBlack, button.text)
				} else {
					tbprint( w / 2 - len(button.text) / 2, h / 2 + 3 + i * 2, termbox.ColorWhite,termbox.ColorBlack, button.text)
				}

			}

		}
		menu.lock.Unlock()
		termbox.Flush()
	}
}