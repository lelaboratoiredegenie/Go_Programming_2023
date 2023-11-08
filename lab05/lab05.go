package main

import (
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"time"
	"html/template"
)

// TODO: Please create a struct to include the information of a video

type Data struct {
    Title string
    Id  string
	ChannelTitle string
	LikeCount string
	ViewCount string
	PublishedAt string
	CommentCount string
}

func with_comma(input string) string {
	var result string

	for i, r := range input {
		if i > 0 && i%3 == 0 {
			result = "," + result
		}
		result = string(r) + result
	}
	return result
}

func YouTubePage(w http.ResponseWriter, r *http.Request) {
	// TODO: Get API token from .env file
	godotenv.Load()
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	// TODO: Get video ID from URL query `v`
	id := r.URL.Query().Get("v")
	// TODO: Get video information from YouTube API
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&key=%s&part=snippet,statistics", id, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		http.ServeFile(w, r, "error.html")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	    http.ServeFile(w, r, "error.html")
	    return
	}
	// TODO: Parse the JSON response and store the information into a struct
	var m map[string]interface{}
	err = json.Unmarshal(body, &m);

	items:= m["items"].([]interface{})
	    
	if len(items) > 0 {
	    item := items[0].(map[string]interface{})

	    snippet := item["snippet"].(map[string]interface{})
	    title := snippet["title"].(string)
		inputdate := snippet["publishedAt"].(string)
		channel := snippet["channelTitle"].(string)

		parsedTime, err := time.Parse(time.RFC3339, inputdate)
		if err != nil {
			http.ServeFile(w, r, "error.html")
			return
		}
		publish := parsedTime.Format("2006年01月02日")

		statistic := item["statistics"].(map[string]interface{})
		like := with_comma(statistic["likeCount"].(string))
		view := with_comma(statistic["viewCount"].(string))
		comment := with_comma(statistic["commentCount"].(string))

		data := Data{
			Title: title,
			Id: id,
			ChannelTitle: channel,
			LikeCount: like,
			ViewCount: view,
			PublishedAt: publish,
			CommentCount: comment,
		}

		err = template.Must(template.ParseFiles("index.html")).Execute(w, data)

    	if err != nil {
    	    http.ServeFile(w, r, "error.html")
    	    return
    	}

	    
	} else {
		http.ServeFile(w, r, "error.html")
		return
	}
	
	// fmt.Println(title)
	// TODO: Display the information in an HTML page through `template`
}

func main() {
	http.HandleFunc("/", YouTubePage)
	log.Fatal(http.ListenAndServe(":8085", nil))
}
