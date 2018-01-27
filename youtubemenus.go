package main

import (
	. "github.com/SimplySerenity/youtube"
	"log"
)

func YoutubeMenu(){
	ClearMenu()
	go InputMenu("Enter a youtube link that you want to watch e.g. https://youtu.be/jAXioRNYy4s",
		func() {

			y := NewYoutube(true)
			err := y.DecodeURL(inputText.Load())

			if err != nil {
				log.Println(err)
				go ClearAndAppend(NewMenuItem("Sorry, but getting the details from YouTube failed :( Maybe try a different link (click to return)", YoutubeMenu))
				return
			}

			//this should fix the bug, I think webm is the only type that youtube serves that omxplayer doesn't support
			for _, stream := range y.StreamList {
				if stream["type"] != "video/webm" {
					PlayLink(stream["url"])
					break
				}
			}

			YoutubeMenu()

		}, MainMenu)
}
