package main

import (
	"github.com/nsf/termbox-go"
	//"log"
)

/*
	The menu, it currently just contains an array of menu items and int of the selected position but will likely expand in the future
*/

type Menu struct {
	Items []MenuItem
	currentlySelected int
}

//gets you a new menu
func NewMenu(items []MenuItem)(Menu){
	return Menu{items, 0}
}

//shift the menu cursor down
func (menu *Menu)CursorDown(){
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

	chunkRadius := (h - yOffset) / 4



	tbprint(w / 2 - len(menu.Items[menu.currentlySelected].text) / 2, h / 2 + yOffset / 2, termbox.ColorBlue, termbox.ColorBlack, menu.Items[menu.currentlySelected].text)

	lowerChunk := menu.Items[max(0, menu.currentlySelected - chunkRadius) : menu.currentlySelected]

	i := 0

	for item := range reverse(lowerChunk){
		tbprint(w / 2 - len(item.text) / 2, h / 2 + yOffset / 2 - (i + 1) * 2, termbox.ColorWhite, termbox.ColorBlack, item.text)
		i++
	}

	upperChunk := menu.Items[menu.currentlySelected + 1 : min(menu.currentlySelected + chunkRadius, len(menu.Items))]

	for i := len(upperChunk) - 1; i >= 0; i-- {
		item := upperChunk[i]
		tbprint(w / 2 - len(item.text) / 2, h / 2 + yOffset / 2 + (i + 1) * 2, termbox.ColorWhite, termbox.ColorBlack, item.text)
	}


}

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