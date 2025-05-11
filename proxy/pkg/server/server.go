package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mmarci96/kubernetes-controller-go/proxy/pkg/watcher"
)

type SpaHandler struct {
	StaticDir   string
	RoutePrefix string
}

func Run() error {
	gameApp := SpaHandler{StaticDir: "./static/game/dist", RoutePrefix: "/game"}
	homeApp := SpaHandler{StaticDir: "./static/client/dist"}
	http.Handle("/game/", gameApp)
	http.Handle("/", homeApp)
	hostUrl := "0.0.0.0:8000"
	fmt.Println("Starting server on: ", hostUrl)
	if err := http.ListenAndServe(hostUrl, nil); err != nil {
		return fmt.Errorf("error starting server: %v", err)
	}
	return nil
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backendServices := watcher.GetBackEndServices()
	fmt.Print("Backend services list logged on serving static files", backendServices)
	path := strings.TrimPrefix(r.URL.Path, h.RoutePrefix)
	fs := http.Dir(h.StaticDir)
	fileServer := http.FileServer(fs)

	f, err := fs.Open(path)
	if err != nil {
		http.ServeFile(w, r, filepath.Join(h.StaticDir, "index.html"))
		return
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			fmt.Printf("Error closing file: %v\n", cerr)
		}
	}()

	stat, err := f.Stat()
	if err != nil || stat.IsDir() {
		http.ServeFile(w, r, filepath.Join(h.StaticDir, "index.html"))
		return
	}

	r.URL.Path = path
	fileServer.ServeHTTP(w, r)
}
