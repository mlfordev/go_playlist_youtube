package main

import (
	"html/template"
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("cmd/webserver/templates/index.html"))
	tpl.Execute(w, nil)
}

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("cmd/webserver/assets"))
    mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", indexHandler)
	fmt.Println("Listening port :3000")
	http.ListenAndServe(":3000", mux)
}
