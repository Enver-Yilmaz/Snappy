package main

/*
	A menu item, it consists of just a string(the text rendered) and a callback for if it's clicked
*/

type MenuItem struct {
	text string
	callback func()
}

//returns a new menu item
func NewMenuItem(text string, callback func())(MenuItem){
	return MenuItem{
		text, callback,
	}
}