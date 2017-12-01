package main

import (
	"github.com/nsf/termbox-go"
	//"log"
	"log"
)

/*
	The menu, it currently just contains an array of menu items and int of the selected position but will likely expand in the future
*/

type Menu struct {
	Items []MenuItem
	currentlySelected int
	inputMode bool
}

//gets you a new menu
func NewMenu(items []MenuItem)(Menu){
	return Menu{items, 0, false}
}

//shift the menu cursor down
func (menu *Menu)CursorDown(){
	log.Println("curDown called!")
	if menu.currentlySelected == len(menu.Items) - 1{
		menu.currentlySelected = 0
	} else {
		menu.currentlySelected += 1
	}
}

//shift the menu cursor up
func (menu *Menu)CursorUp(){
	if menu.currentlySelected == 0{
		menu.currentlySelected = len(menu.Items) - 1
	} else {
		menu.currentlySelected -= 1
	}
}

//trigger the callback of the selected item
func (menu *Menu)MenuSelect(){
	menu.Items[menu.currentlySelected].callback()
}

//draw the menu
//the scrolling is honestly a really dumb solution but I'm a really dumb guy
func (menu *Menu)DrawMenu(yOffset int){
	w, h := termbox.Size()


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
		text := "Input text"
		tbprint(w / 2 - len(text) / 2, h / 2 - 5, termbox.ColorWhite, termbox.ColorBlack, text)
		tbprint(w / 2 - len(inputText) / 2, h / 2, termbox.ColorWhite, termbox.ColorBlack, inputText)
		tbprint(w / 2 - len("Input: ") - len(inputText) / 2, h / 2, termbox.ColorWhite, termbox.ColorBlack, "Input: ")

		for i, button := range menu.Items{

			if i == menu.currentlySelected{
				tbprint( w / 2 - len(button.text) / 2, h / 2 + 3 + i * 2, termbox.ColorBlue, termbox.ColorBlack, button.text)
			} else {
				tbprint( w / 2 - len(button.text) / 2, h / 2 + 3 + i * 2, termbox.ColorWhite,termbox.ColorBlack, button.text)
			}

		}

	}


}

//a weird way to reverse, thanks stackoverflow
func reverse(lst []MenuItem) chan MenuItem {
	ret := make(chan MenuItem)
	go func() {
		for i := range lst {
			ret <- lst[len(lst)-1-i]
		}
		close(ret)
	}()
	return ret
}

// because go lacks ternary
func min(a, b int)int{
	if a <= b {
		return a
	}
	return b
}

func max(a, b int)int{
	if a >= b{
		return a
	}
	return b
}