package main

import (
	"embed"
	"encoding/json"
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
		http.FileServer(http.FS(subFS)).ServeHTTP(w, r)
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
