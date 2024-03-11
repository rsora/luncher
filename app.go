package main

import (
	"bytes"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

const suggestionTemplate = `<html>
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
	<div class=rotated><h1>{{ .Recipe }}</h1></div>
</body>
</html>`

const recipeTemplate = `<html>
<head>
<style>
body {
  font-family: 'Courier New', monospace;
}
</style>
</head>
<body>
	<div>
	<h1><a href="status">Ciambellone per tutte le ore</a></h1>
	<h2>Ingredienti</h2>
	<ul>
		<li>300 gr. farina 00</li>
		<li>250 gr. zucchero semolato</li>
		<li>200 gr. burro</li>
		<li>4 uova  </li>
		<li>13 gr. cremor tartaro </li>
		<li>mandorle qb</li>
		<li>gocce di cioccolato qb</li>
	</ul>
	<h2>Preparazione</h2>
	<ol>
		<li>Accendere il forno e tostare mandorle nel forno mentre si scalda</li>
		<li>Scaldare burro nel microonde 10 secondi per ammorbidirlo</li>
		<li>Mescolare burro e zucchero nella planetaria</li>
		<li>Aggiungere le uova ad una ad una per farle incoporare bene</li>
		<li>mescolare in una ciotola le farine, il cremor tartaro</li>
		<li>Aggiungere le farine a mano a mano, appena scompaiono, spegnere la planetaria</li>
	</ol>
	<h2>Cottura</h2>
	Forno sopra e sotto 170 C per 40 minuti
	<h2>Note</h2>
	<ul>
		<li>imburrare la teglia quadrata piccola di acciaio</li>
	</ul>
	</div>
</body>
</html>`

type Suggestion struct {
	Recipe string
}

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/status", a.getStatus).Methods("GET")
	a.Router.HandleFunc("/daily", a.getSuggestion).Methods("GET")
	a.Router.HandleFunc("/recipes", a.getRecipes).Methods("GET")

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) getSuggestion(w http.ResponseWriter, r *http.Request) {
	SuggestionList := []Suggestion{
		{Recipe: "Gazpacho senza cetrioli con polpa di peperoni e crostini di pane ai cereali"},
		{Recipe: "Spaghetti grossi con sugo rosso tonno e olive"},
		{Recipe: "Insalata di pollo"},
		{Recipe: "Insalata di riso"},
		{Recipe: "Insalata tonno, pomodoro e avocado"},
		{Recipe: "Pasta con ragu'"},
		{Recipe: "Frittata con le zucchine"},
		{Recipe: "Lenticchie rosse e salsiccia in umido"},
		{Recipe: "Crescia sfogliata fatta in casa con pomodoro mozzarella crudo e guacamole"},
		{Recipe: "Spezzatino con sedano e puree di patate e zucchine"},
	}
	s := time.Now().Unix() % int64(len(SuggestionList))
	log.Println(s)
	sp := AddSimpleTemplate(suggestionTemplate, SuggestionList[s])
	respondWithHTML(w, http.StatusOK, sp)
}

func (a *App) getRecipes(w http.ResponseWriter, r *http.Request) {

	sp := AddSimpleTemplate(recipeTemplate, Suggestion{})
	respondWithHTML(w, http.StatusOK, sp)
}

func (a *App) getStatus(w http.ResponseWriter, r *http.Request) {
	respondWithHTML(w, http.StatusOK, "ok")
}

func AddSimpleTemplate(a string, b Suggestion) string {
	tmpl := template.Must(template.New("suggestion.recipe").Parse(a))
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, b)
	if err != nil {
		panic(err)
	}
	s := buf.String()
	return s
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithHTML(w, code, message)
}

func respondWithHTML(w http.ResponseWriter, code int, payload string) {

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	w.Write([]byte(payload))
}
