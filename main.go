package main

import (
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	startApi()
}

func startApi() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Get("/hello", helloHandler)
	})

	r.Handle("/*", http.HandlerFunc(frontendHandler))

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	listenUrl := host + ":" + port

	log.Printf("Starting server on %s", listenUrl)
	err := http.ListenAndServe(listenUrl, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

//go:embed all:frontend/dist
var distFS embed.FS

func frontendHandler(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("ENV")
	viteUrl := os.Getenv("VITE_URL")

	if viteUrl == "" {
		viteUrl = "http://localhost:5173"
	}

	if env == "dev" {
		target, err := url.Parse(viteUrl)
		if err != nil {
			log.Fatalf("Could not parse VITE_URL: %s", viteUrl)
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	} else {
		subFS, err := fs.Sub(distFS, "frontend/dist")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the request is for a file (has an extension) or if the file exists
		path := r.URL.Path

		// Try to open the file
		file, err := subFS.Open(path[1:]) // Remove leading slash
		if err == nil {
			file.Close()
			// File exists, serve it normally
			http.FileServer(http.FS(subFS)).ServeHTTP(w, r)
		} else {
			// File doesn't exist, serve index.html for client-side routing
			// This enables SPA functionality
			indexFile, err := subFS.Open("index.html")
			if err != nil {
				http.Error(w, "index.html not found", http.StatusNotFound)
				return
			}
			defer indexFile.Close()

			// Serve index.html content
			stat, _ := indexFile.Stat()
			http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
		}
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"msg": "Hello, Gopher!"}

	jsonResponse, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
