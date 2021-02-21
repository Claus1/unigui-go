package unigui

import (	
	"flag"
	"log"
	"net/http"
	"strings"
)


var( 
	addr = flag.String("addr", ":1234", "websocket service address")
	upload_dir = "upload"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	
	path := r.URL.Path
	i := strings.Index(path, "?")

	if i != -1{
		path = path[:i]
	}
	path = strings.ReplaceAll(path, "%20"," ")
	uplPath := F("/%s/",path)
	if strings.Index(path, uplPath) != 0{
		path = "web" + path
	}
		
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
		http.ListenAndServe(":8000", mxHTTP)
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}	
}
