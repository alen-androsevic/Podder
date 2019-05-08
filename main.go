package main

import (
	"github.com/SlyMarbo/rss"
	"bufio"
	"log"
	"os"
        "html/template"
	"net/http"
)

var tpl *template.Template

type Podcast struct {
	Title string
	Link string
}

var pods []Podcast

type ViewData struct{
    Podcasts []Podcast
}

func main() {
	feeds, err := os.Open("feeds")
	if err != nil {
		log.Fatal(err)
	}

	defer feeds.Close()
	scanner := bufio.NewScanner(feeds)
	for scanner.Scan() {
		feed, err := rss.Fetch(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		for index := 0; index < len(feed.Items); index++ {
			url := feed.Items[index].Enclosures[0].URL
			podcast := Podcast{
				Title: feed.Items[index].Title,
				Link: url,
			}
			pods = append(pods, podcast)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	tpl, err = template.ParseFiles("feeds.gohtml")
	if err != nil {
          log.Fatal(err)
        }

	http.HandleFunc("/", handler)
	http.HandleFunc("/style.css", serveResource)
	http.HandleFunc("/favicon-32x32.png", serveResource)
	http.HandleFunc("/favicon-16x16.png", serveResource)
	http.HandleFunc("/site.webmanifest", serveResource)
	http.ListenAndServe(":80", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	vd := ViewData{Podcasts: pods,}

	err := tpl.Execute(w, vd)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "." + req.URL.Path
	http.ServeFile(w, req, path)
}
