package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/mmarci96/kubernetes-controller-go/proxy/pkg/watcher"
)

var proxyClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
}

type SpaHandler struct {
	StaticDir   string
	RoutePrefix string
}

func Run() error {
	http.HandleFunc("/socket.io/", handler)

	gameApp := SpaHandler{StaticDir: "./static/game/dist", RoutePrefix: "/game"}
	homeApp := SpaHandler{StaticDir: "./static/client/dist"}
	http.Handle("/game/", gameApp)
	http.Handle("/", homeApp)
	hostUrl := "0.0.0.0:80"
	fmt.Println("Starting server on: ", hostUrl)
	if err := http.ListenAndServe(hostUrl, nil); err != nil {
		return fmt.Errorf("error starting server: %v", err)
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	gameId := r.URL.Query().Get("gameId")
	playerId := r.URL.Query().Get("playerId")
	backends := watcher.GetGameServers()
	if len(backends) == 0 {
		http.Error(w, "No backends available", http.StatusServiceUnavailable)
		return
	}
	// Implement your load balancing logic here using 'backends'
	// Example: Round Robin, Random selection, etc.
	fmt.Println("Backend on handler", backends)
	fmt.Println("Gameid on request", gameId)
	fmt.Println("Playerid on request", playerId)
	backend := backends[0]
	proxyRequest(backend, w, r)

}

func proxyRequest(backend string, w http.ResponseWriter, r *http.Request) {
	// Create target URL
	target := &url.URL{
		Scheme: "http",            // Assuming HTTP backend
		Host:   backend + ":8080", // Use your actual service port
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = proxyClient.Transport
	// Modify the request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Remove gameId and playerId from query parameters
		query := req.URL.Query()
		query.Del("gameId")
		query.Del("playerId")
		req.URL.RawQuery = query.Encode()

		// Preserve original path and set correct host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		// Optional: Add custom headers
		req.Header.Set("X-Forwarded-For", r.RemoteAddr)
	}

	// Custom error handling
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[ERROR] Proxy error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
	}

	// Serve the request using our custom client
	proxy.ServeHTTP(w, r)
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
