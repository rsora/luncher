package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/status", a.getStatus).Methods("GET")
	a.Router.HandleFunc("/daily", a.getProduct).Methods("GET")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	respondWithHTML(w, http.StatusOK, nil)
}

func (a *App) getStatus(w http.ResponseWriter, r *http.Request) {
	respondWithHTML(w, http.StatusOK, nil)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithHTML(w, code, map[string]string{"error": message})
}

func respondWithHTML(w http.ResponseWriter, code int, payload interface{}) {

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	w.Write([]byte("ok"))
}
