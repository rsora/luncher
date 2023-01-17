package handler

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"text/template"
	"time"
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
	<div class=rotated><h1>{{ .Recipe }}</h1></div>
</body>
</html>`

type Suggestion struct {
	Recipe string
}

func respondWithHTML(w http.ResponseWriter, code int, payload string) {

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	w.Write([]byte(payload))
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

func GetSuggestion(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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
	sp := AddSimpleTemplate(htmlTemplate, SuggestionList[s])
	respondWithHTML(w, http.StatusOK, sp)
	return nil
}
