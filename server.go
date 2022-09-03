package unigui

import (
	"bytes"
	"fmt"
	"io"
	"log"	
	"io/ioutil"
	h "net/http"
	"os"
	"strings"	
	"github.com/hashicorp/go-getter"
)

var (
	ResourcePort = ":8000"
	WsocketPort  = ":8000"
	UploadDir    = "upload"
	SocketIp     = "localhost"
	AppName      = "Unigui" 
	config    map[string]string
	funcServeMain func(w h.ResponseWriter, r *h.Request)	
)

func ReadConfig() {
	file, err := os.Open("config")
	if err != nil {
		Print(err)
	}
	defer file.Close()
	byteValue, _ := ioutil. ReadAll(file)	
	config = make(map[string]string)
	for _, str := range strings.Split(string(byteValue), "\n") {
		str = strings.TrimSpace(str)
		if str != "" {
			arr := strings.Split(str, "=")
			if len(arr) == 2 {
				config[strings.TrimSpace(arr[0])] = strings.TrimSpace(arr[1])
			}
		}
	}
}


func GetConfig(param string) string {
	return config[param]
}

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
	fpath := F("./web%s", strings.ReplaceAll(path, "%20", " "))

	if r.Method == "GET" {
		if path == "/" {
			//serveMain(w, r, fpath)
			h.ServeFile(w, r, fpath + "index.html")
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

	ReadConfig()

	ResourcePort = GetConfig("port")
	AppName = GetConfig("appname")
	
	hub := newHub()
	go hub.run()

	h.HandleFunc("/ws", func(w h.ResponseWriter, r *h.Request) {
		serveWs(hub, w, r)
	})	
	
	h.HandleFunc("/", serveHttp)	
		
	err := h.ListenAndServe(ResourcePort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
