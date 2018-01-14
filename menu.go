package main

import (
	"sync"
)

/*
	The menu, it's really just a slice of menu items that allow packed with some built in methods to make it thread safe
*/


type Menu struct {
	Items []MenuItem
	currentlySelected int
	inputMode bool
	lock *sync.Mutex
}

//gets you a new menu
//should only be called once
func NewMenu(items []MenuItem)(Menu){
	return Menu{items, 0, false, new(sync.Mutex)}
}

//appends items from the append channel
func AppendMenu(items ...MenuItem){
	menu.lock.Lock()
	for _, item := range items{
		menu.Items = append(menu.Items, item)
	}
	menu.lock.Unlock()
	drawChan <- 1
}

func SetInputMode(inputMode bool){
	menu.lock.Lock()
	menu.inputMode = inputMode
	menu.lock.Unlock()
}

//for building menus
//because go doesn't have overloads or optional arguments
func ClearAndAppend(items ...MenuItem){
	ClearMenu() //yeah we're going to block
	menu.lock.Lock()
	for _, item := range items{
		menu.Items = append(menu.Items, item)
	}
	menu.lock.Unlock()
	drawChan <- 1
}

//clears the menu
func ClearMenu(){
	menu.lock.Lock()
	menu.Items = []MenuItem{}
	menu.currentlySelected = 0
	menu.lock.Unlock()
	drawChan <- 1
}

//shift the menu cursor down
func CursorDown(){
	menu.lock.Lock()
	if menu.currentlySelected == len(menu.Items) - 1 {
		menu.currentlySelected = 0
	} else {
		menu.currentlySelected += 1
	}
	menu.lock.Unlock()
	drawChan <- 1
}

//shift the menu cursor up
func CursorUp(){
	menu.lock.Lock()
	if menu.currentlySelected == 0 {
		menu.currentlySelected = len(menu.Items) - 1
	} else {
		menu.currentlySelected -= 1
	}
	menu.lock.Unlock()
	drawChan <- 1
}

//trigger the callback of the selected item
func MenuSelect(){
	go menu.Items[menu.currentlySelected].callback()
	drawChan <- 1
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