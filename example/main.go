package main

import (
	"html/template"
	"log"
	"net/http"

	"viewer"
	"viewer/example/state"
)

var t = template.Must(template.New("").Funcs(viewer.Funcs()).ParseFiles(
	"index.html",
	"hello.tmpl",
))

func index(w http.ResponseWriter, r *http.Request) {
	data := struct{ viewer.Templater }{
		viewer.Templater{
			T: t,
			InitialState: state.State{
				Name: "Damien",
			},
		},
	}

	if err := t.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Println("rendering error: " + err.Error())
	}
}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/static/", http.FileServer(http.Dir(".")))

	log.Println("listening on 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
