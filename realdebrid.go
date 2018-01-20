package main

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"github.com/buger/jsonparser"
	"log"
	"github.com/nsf/termbox-go"
)

//uses real-debrid to unrestrict a link that we likely found with Alluc
func unrestrict(link string, title string){
	go func(){

		form := url.Values{}
		form.Add("link", link)

		resp, err := http.PostForm("https://api.real-debrid.com/rest/1.0/unrestrict/link?auth_token="+realDebridKey, form)

		if err != nil {
			log.Println("failed to connect with real-debrid API")
			return
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			termbox.Close()
			log.Fatal(err)
		}
		streamURL, err := jsonparser.GetString(body, "download")

		if err != nil {
			log.Println("couldn't debrid unrestrict this link: " + link)
			log.Println("Error body:" + string(body))
		} else {
			log.Println(streamURL)

			go AppendMenu(NewMenuItem(title, func() {
				PlayLink(streamURL)
				inPlayback.Store(true)
			}))
		}

	}()
}