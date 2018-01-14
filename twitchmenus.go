package main

import (
	"net/http"
	"log"
	//"math/rand"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"fmt"
	//"strconv"
	"github.com/grafov/m3u8"
	"os/exec"
	//"net/url"
	"net/url"
	"github.com/nsf/termbox-go"
)

const twitchClientID = "dlpf1993tub698zw0ic6jlddt9e893"
const tokenApi = "https://api.twitch.tv/api/channels/%s/access_token?client_id=%s"
const playlistApi = "https://usher.ttvnw.net/api/channel/hls/%s.m3u8?token=%s&sig=%s&allow_source=true&player_backend=html5"


func TwitchSearchMenu(){
	go ClearAndAppend(
		NewMenuItem("Enter", func(){
			if inputText.Load() != "" {
				TwitchStreamsMenu(inputText.Load())
				go SetInputMode(false) //i'm actually not sure if there is any reason to spawn another goroutine for this TODO
			}
		}),
		NewMenuItem("Back", func(){
			MainMenu()
			go SetInputMode(false)
		}),
	)
	go SetInputMode(true)
}

func TwitchStreamsMenu(channel string){
	resp, err := http.Get(fmt.Sprintf(tokenApi, channel, twitchClientID))

	if err != nil {
		go ClearAndAppend(
			NewMenuItem("The Twitch api request failed, maybe it's a connectivity problem? (click to return)", func(){
				TwitchSearchMenu()
			}))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	resp.Body.Close()

	if err != nil {
		termbox.Close() //maybe so things don't wonk out
		panic(err)	//if we really couldn't read the body, then something went real wrong
	}

	token, err := jsonparser.GetString(body, "token")

	if err != nil {
		go ClearAndAppend(
			NewMenuItem("Twitch didn't authenticate us :( file an issue because it's possible the API changed (click to return)", func(){
				TwitchSearchMenu()
			}))
		return
	}

	sig, err := jsonparser.GetString(body, "sig")

	if err != nil {	//really should have failed at the previous at this point
		go ClearAndAppend(
			NewMenuItem("Twitch didn't authenticate us :( file an issue because it's possible the API changed (click to return)", func(){
				TwitchSearchMenu()
			}))
		return
	}

	resp, err = http.Get(fmt.Sprintf(playlistApi, channel, url.QueryEscape(token), sig))

	if err != nil {
		go ClearAndAppend(
			NewMenuItem("The Twitch api request failed, maybe it's a connectivity problem? (click to return)", func(){
				TwitchSearchMenu()
			}))
		return
	}

	p, listType, err := m3u8.DecodeFrom(resp.Body, false)

	TBcheck(err)

	if listType != m3u8.MASTER {
		go ClearAndAppend(
			NewMenuItem("couldn't find a stream by that name, maybe it's offline? (click to return)", func() {
				TwitchSearchMenu()
			}))
		return
	}
	ClearMenu()

	captureURI := func(uri string) func() {	//so we can use it while in a range loop of the playlist
		return func(){
			if !desktop {
				command := exec.Command("omxplayer", "-b", "-o", "hdmi", "--live", uri)
				log.Println(uri)
				err = command.Run()
			} else {
				command := exec.Command("vlc", uri)
				log.Println(uri)
				err = command.Run()
			}
			inPlayback.Store(true)
		}
	}

	for _, variant := range p.(*m3u8.MasterPlaylist).Variants {

		if variant.Video == "chunked"{
			variant.Video = "Source"
		}
		log.Println(variant.URI)
		log.Println(variant.Video)
		AppendMenu(NewMenuItem(variant.Video, captureURI(variant.URI))) //TODO consider using this style more, I kinda like it
	}
	AppendMenu(NewMenuItem("Return", func(){
		TwitchSearchMenu()
	}))
}