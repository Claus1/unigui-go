package unigui

import (
	"flag"
	"fmt"
	"io"
	"log"
	h "net/http" 
	"os"
	"strings"
	"github.com/hashicorp/go-getter"
	"bytes"		
)

var (
	ResourcePort = ":8000"
	WsocketPort  = ":1234"
	UploadDir    = "upload"
	SocketIp     = "localhost"
	mainJs = "/main.dart.js"
	funcServeMain func(w h.ResponseWriter, r *h.Request)
)

//download web files if do not exist
func init() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Print(dir)
	}
	dir += "/web"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Print("downloading web files...")
		getter.Get(dir, "github.com/Claus1/unigui-go//web")
		fmt.Print(" done!. unigui is ready to use.")
	}
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", h.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", h.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", h.StatusInternalServerError
}

func serveMain(w h.ResponseWriter, r *h.Request, fpath string) {
	if funcServeMain == nil {
		f, err := os.Open(fpath)
		if err != nil {
			msg, code := toHTTPError(err)
			h.Error(w, msg, code)
			return
		}
		defer f.Close()	

		fileInfo, _ := f.Stat()
		var size int64 = fileInfo.Size()
	
		buffer := make([]byte, size)
		
		f.Read(buffer)

		mainBuffer := bytes.ReplaceAll(buffer, []byte("localhost"), []byte(SocketIp))
	
		fileBytes := bytes.NewReader(mainBuffer) // converted to io.ReadSeeker type	

		funcServeMain = func(w h.ResponseWriter, r *h.Request){
			h.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), fileBytes)
		}
	}
	funcServeMain(w, r)	
}

func serveHttp(w h.ResponseWriter, r *h.Request) {

	path := r.URL.Path
	i := strings.Index(path, "?")

	if i != -1 {
		path = path[:i]
	}
	fpath := F("web%s", strings.ReplaceAll(path, "%20", " "))

	if r.Method == "GET" {
		if path == mainJs && SocketIp != "localhost"{
			serveMain(w, r, fpath)
		} else{
			h.ServeFile(w, r, fpath)
		}

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
						w.WriteHeader(h.StatusInternalServerError)
						panic(err)
					}
				} else {
					w.WriteHeader(h.StatusInternalServerError)
					panic(err)
				}
			}
		}
	} else {
		h.Error(w, "Method not allowed", h.StatusMethodNotAllowed)
	}
}

func Start() {

	flag.Parse()
	hub := newHub()
	go hub.run()

	mxHTTP := h.NewServeMux()
	mxHTTP.HandleFunc("/", serveHttp)	
	go func() {
		h.ListenAndServe(ResourcePort, mxHTTP)
	}()

	h.HandleFunc("/", func(w h.ResponseWriter, r *h.Request) {
		serveWs(hub, w, r)
	})
	err := h.ListenAndServe(WsocketPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
