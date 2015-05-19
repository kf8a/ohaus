package main

import (
	// "encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type connection struct {
	ws   *websocket.Conn
	send chan []byte
	d    *dataSource
}

func (c *connection) reader() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			return
		}
	}
	c.ws.Close()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ScaleHandler(instrument *dataSource, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte), ws: ws, d: instrument}
	c.d.register <- c
	defer func() { c.d.unregister <- c }()
	c.reader()
}

func StartRecordingHandler(d *dataSource, w http.ResponseWriter, r *http.Request) {
}

func main() {
	var test bool
	flag.BoolVar(&test, "test", false, "use a random number generator instead of a live feed")
	flag.Parse()

	instrument := newDataSource()
	go instrument.read(test)

	r := mux.NewRouter()

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ScaleHandler(instrument, w, r)
	})

	r.HandleFunc("/record", func(w http.ResponseWriter, r *http.Request) {
		StartRecordingHandler(instrument, w, r)
	})

	http.Handle("/", r)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.ListenAndServe("127.0.0.1:8081", nil)
}
