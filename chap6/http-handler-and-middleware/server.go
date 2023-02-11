package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type appConfig struct {
	logger *log.Logger
}

type app struct {
	config  appConfig
	handler func(
		w http.ResponseWriter, r *http.Request, config appConfig,
	)
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	a.handler(w, r, a.config)
	a.config.logger.Printf(
		"path=%s method=%s duration=%f", r.URL.Path, r.Method,
		time.Now().Sub(startTime).Seconds(),
	)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	config.logger.Println("Handling healthcheck request")
	fmt.Fprintf(w, "ok")
}

func apiHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	config.logger.Println("Handling API request")
	fmt.Fprintf(w, "Hello, World!")
}

func panicHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	panic("I panicked")
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			h.ServeHTTP(w, r)
			log.Printf(
				"protocol=%s path=%s method=%s duration=%f",
				r.Proto,
				r.URL.Path,
				r.Method,
				time.Now().Sub(startTime).Seconds(),
			)
		})
}

func panicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rValue := recover(); rValue != nil {
				log.Println("panic detected when handling request", rValue)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Unexpected server error occured")
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func setupHandlers(mux *http.ServeMux, config appConfig) {
	mux.Handle("/healthz", &app{config: config, handler: healthCheckHandler})
	mux.Handle("/api", &app{config: config, handler: apiHandler})
	mux.Handle("/panic", &app{config: config, handler: panicHandler})
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}
	config := appConfig{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	mux := http.NewServeMux()
	setupHandlers(mux, config)
	m := loggingMiddleware(panicMiddleware(mux))
	log.Fatal(http.ListenAndServe(listenAddr, m))
}
