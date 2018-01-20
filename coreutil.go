package main

import "os/exec"

func PlayLink(link string){
	if !desktop {
		command := exec.Command("omxplayer", "-b", "-o", "hdmi", "--threshold", "0", link)
		command.Run()
	} else {
		command := exec.Command("vlc", link)
		command.Run()
	}
}