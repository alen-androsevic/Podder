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

type Category struct {
	Title string
	Podcasts []Podcast
}

type Podcast struct {
	Title string
	Link string
}

var cats []Category

type ViewData struct{
    Categories []Category
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

		var pods []Podcast
		for index := 0; index < len(feed.Items); index++ {
			url := feed.Items[index].Enclosures[0].URL
			podcast := Podcast{
				Title: feed.Items[index].Title,
				Link: url,
			}
			pods = append(pods, podcast)
		}
		category := Category{
			Title: feed.Title,
			Podcasts: pods,
		}
		cats = append(cats, category)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	tpl, err = template.ParseFiles("feeds.gohtml")
	if err != nil {
          log.Fatal(err)
        }

	http.HandleFunc("/", handler)
	http.HandleFunc("/app.js", serveResource)
	http.HandleFunc("/style.css", serveResource)
	http.HandleFunc("/favicon-32x32.png", serveResource)
	http.HandleFunc("/favicon-16x16.png", serveResource)
	http.HandleFunc("/site.webmanifest", serveResource)

	port := ":" + os.Getenv("PORT")

	if port == ":" {
		port = ":3000"
	}

	http.ListenAndServe(port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	vd := ViewData{Categories: cats,}

	err := tpl.Execute(w, vd)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "." + req.URL.Path
	http.ServeFile(w, req, path)
}
