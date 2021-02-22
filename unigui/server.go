package unigui

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

var (
	ResourcePort = ":8000"
	WsocketPort  = ":1234"
	UploadDir    = "upload"
)

func serveHome(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	i := strings.Index(path, "?")

	if i != -1 {
		path = path[:i]
	}	
	path = F("web/%s/", strings.ReplaceAll(path, "%20", " "))
		
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, path)
}

func Start() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	mxHTTP := http.NewServeMux()
	mxHTTP.HandleFunc("/", serveHome)
	go func() {
		http.ListenAndServe(ResourcePort, mxHTTP)
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(WsocketPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
