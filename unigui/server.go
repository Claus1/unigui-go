package unigui

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

	if r.Method == "GET" {
		http.ServeFile(w, r, path)
	} else if r.Method == "POST" {
		err := r.ParseMultipartForm(10 << 20) // grab the multipart form
		if err != nil {
			fmt.Fprintln(w, err)
		}
		for _, fheaders := range r.MultipartForm.File {
			for _, hdr := range fheaders {
				// open uploaded
				infile, _ := hdr.Open()
				// open destination
				var outfile *os.File
				if outfile, err = os.Create(F("web/%s/%s", UploadDir, hdr.Filename)); err == nil {
					defer outfile.Close()
					if _, err := io.Copy(outfile, infile); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						panic(err)
					}
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					panic(err)
				}
			}
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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
