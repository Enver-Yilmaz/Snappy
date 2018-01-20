package main

import (
	"net/http"
	"log"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"fmt"
	"github.com/grafov/m3u8"
	"net/url"
	"github.com/nsf/termbox-go"
	"strings"
)

const twitchClientID = "dlpf1993tub698zw0ic6jlddt9e893"
const tokenApi = "https://api.twitch.tv/api/channels/%s/access_token?client_id=%s"
const playlistApi = "https://usher.ttvnw.net/api/channel/hls/%s.m3u8?token=%s&sig=%s&allow_source=true&player_backend=html5"

//Let's the user input the channel they want to search for on Twitch
func TwitchSearchMenu(){
	go InputMenu("Enter the twitch channel you want to watch e.g. vinesauce", func(){
		TwitchStreamsMenu(strings.ToLower(inputText.Load()))
	}, func(){
		SearchMenu()
	})
}

//uses the unoffical twitch endpoints to find and parse the m3u8 file that twitch uses for its livestreams
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
			PlayLink(uri)
			inPlayback.Store(true)
		}
	}

	for _, variant := range p.(*m3u8.MasterPlaylist).Variants {

		if variant.Video == "chunked" { 	//following twitch styling
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