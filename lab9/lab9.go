package main

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"strings"

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

func InitObservable() {
	swearWords := loadWordsFromFile("swear_word.txt")
	sensitiveNames := loadSensitiveNames("sensitive_name.txt")

	ObservableMsg = rxgo.FromChannel(messages).
		Filter(func(item interface{}) bool {
			if message, ok := item.(string); ok {
				for swearWord := range swearWords {
					if strings.Contains(message, swearWord) {
						return false
					}
				}
			}
			return true
		}).
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			if message, ok := item.(string); ok {
				for name, modifiedName := range sensitiveNames {
					message = strings.ReplaceAll(message, name, modifiedName)
				}
				return message, nil
			}
			return item, nil
		})
}

func loadWordsFromFile(filename string) map[string]bool {
	words := make(map[string]bool)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open %s: %v", filename, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words[scanner.Text()] = true
	}
	return words
}

func loadSensitiveNames(filename string) map[string]string {
	names := make(map[string]string)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open %s: %v", filename, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := scanner.Text()
		runes := []rune(name)
		if len(runes) >= 2 {
			names[name] = string(runes[0]) + "*" + string(runes[2:])
		} else {
			names[name] = name
		}
	}
	return names
}

func main() {
	InitObservable()
	go broadcaster()
	http.HandleFunc("/wschatroom", wshandle)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("server start at :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}