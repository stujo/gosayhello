package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"html/template"
	"time"
)

var templates = template.Must(template.ParseFiles("index.html"))

type Greeting struct {
	Version string
	Time    string
}

func main() {

	http.HandleFunc("/", sayhello)

	var hostname = "localhost"
	var port = "3337"

	if os.Getenv("HOST") != "" {
		hostname = os.Getenv("HOST")
	}

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	//Using Default localhost
	bind := fmt.Sprintf("%s:%s", hostname, port)

	fmt.Printf("listening on %s...", bind)

	err := http.ListenAndServe(bind, nil)

	if err != nil {
		panic(err)
	}
}

func sayhello(res http.ResponseWriter, req *http.Request) {

	greeting := &Greeting{Version: runtime.Version(), Time: time.Now().String()}

	err := templates.ExecuteTemplate(res, "index.html", greeting)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
