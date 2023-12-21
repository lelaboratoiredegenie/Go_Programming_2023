package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
	"context"
	"regexp"
	"unicode/utf8"
	"github.com/gorilla/websocket"
	"github.com/reactivex/rxgo/v2"
)

type client chan<- string // an outgoing message channel

var (
	entering      = make(chan client)
	leaving       = make(chan client)
	messages      = make(chan rxgo.Item) // all incoming client messages
	ObservableMsg = rxgo.FromChannel(messages)
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	MessageBroadcast := ObservableMsg.Observe()
	for {
		select {
		case msg := <-MessageBroadcast:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg.V.(string)
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func clientWriter(conn *websocket.Conn, ch <-chan string) {
	for msg := range ch {
		conn.WriteMessage(1, []byte(msg))
	}
}

func wshandle(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "你是 " + who + "\n"
	messages <- rxgo.Of(who + " 來到了現場" + "\n")
	entering <- ch

	defer func() {
		log.Println("disconnect !!")
		leaving <- ch
		messages <- rxgo.Of(who + " 離開了" + "\n")
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		messages <- rxgo.Of(who + " 表示: " + string(msg))
	}
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}


func InitObservable() {

	swearWords, err := readLines("swear_word.txt")
	if err != nil {
		log.Fatal(err)
	}

	sensitiveNames, err := readLines("sensitive_name.txt")
	if err != nil {
		log.Fatal(err)
	}

	ObservableMsg = rxgo.FromChannel(messages)

	filterFn := func(item interface{}) bool {
		msg := item.(string)
		for _, word := range swearWords {
			if strings.Contains(msg, word) {
				return false
			}
		}
		return true
	}

	mapFn := func(i interface{}) (interface{}, error) {
		msg := i.(string)
		
		for _, name := range sensitiveNames {
			pattern := regexp.QuoteMeta(name)
        	re := regexp.MustCompile(pattern)

        	replacer := func(match string) string {
            	if utf8.RuneCountInString(name) == 3 {
            	    runes := []rune(match)
            	    runes[1] = []rune("*")[0] 
            	    return string(runes)
            	} else if utf8.RuneCountInString(name) == 2 {
            	    runes := []rune(match)
            	    runes[1] = []rune("*")[0]
            	    return string(runes)
            	}
            	return match
        	}
        	msg = re.ReplaceAllStringFunc(msg, replacer)
    	}
    	return msg, nil
	}

	ObservableMsg = ObservableMsg.
		Filter(func(i interface{}) bool {
			return filterFn(i)
		}).
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			transed, _ := mapFn(item)
			return transed, nil
		})
}

func main() {
	InitObservable()
	go broadcaster()
	http.HandleFunc("/wschatroom", wshandle)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("server start at :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}