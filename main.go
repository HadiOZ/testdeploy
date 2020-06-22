package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", corsmiddleware(indexHandel))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	fmt.Println("server run on port 8085")
	if err := http.ListenAndServe(":8085", nil); err != nil {
		log.Fatal(err)
	}
}

func corsmiddleware(hendel http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		hendel(w, r)
	}
}

func indexHandel(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
	}
	tmp := template.Must(t, err)
	if err != nil {
		http.Error(w, "can't render template", http.StatusInternalServerError)
	}
	tmp.Execute(w, nil)
}
