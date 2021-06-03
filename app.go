package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const htmlTemplate = `<html>
<head>
	<style>
	

	.rotated 
		{
			-webkit-transform: translate(-50%, -50%) rotate(90deg);
			-moz-transform: translate(-50%, -50%) rotate(90deg);
			-ms-transform: translate(-50%, -50%) rotate(90deg);
			-o-transform: translate(-50%, -50%) rotate(90deg);
			position: fixed;
			top: 50%;
			left: 50%;
			text-align: center;
			font-family: monospace; 
			width: 80%;
		}
	</style>
	</head>
	
<body>
	<div class=rotated><h1>Riso alla cantonese saltato con uova, prosciutto, piselli</h1></div>
</body>
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
