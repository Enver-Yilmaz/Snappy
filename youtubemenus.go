package main

import (
	. "github.com/SimplySerenity/youtube"
	"log"
)

func YoutubeMenu(){
	ClearMenu()
	go InputMenu("Enter a youtube link that you want to watch e.g. https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		func() {

			y := NewYoutube(true)
			err := y.DecodeURL(inputText.Load())

			if err != nil {
				log.Println(err)
				go ClearAndAppend(NewMenuItem("Sorry, but getting the details from YouTube failed :( Maybe try a different link (click to return)", YoutubeMenu))
				return
			}

			PlayLink(y.GetDownloadUrl())
			YoutubeMenu()

		}, MainMenu)
}
