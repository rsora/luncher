package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const htmlTemplate = `<html>
<head>
	<style>
	th 
		{
		vertical-align: bottom;
		text-align: center;
		}

	th span 
		{
		-ms-writing-mode: tb-rl;
		-webkit-writing-mode: vertical-rl;
		writing-mode: vertical-rl;
		transform: rotate(180deg);
		white-space: nowrap;
		}
	</style>
	</head>
	
<table>
	<tr>
		<th><span>Bene bene Molto bene</span></th>
	</tr>
</table>
</html>`

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
	respondWithHTML(w, http.StatusOK, htmlTemplate)
}

func (a *App) getStatus(w http.ResponseWriter, r *http.Request) {
	respondWithHTML(w, http.StatusOK, "ok")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithHTML(w, code, message)
}

func respondWithHTML(w http.ResponseWriter, code int, payload string) {

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	w.Write([]byte(payload))
}
