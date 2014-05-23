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

func heroku_binding() (hostname string, port string) {
	hostname = os.Getenv("HOST")
	port = os.Getenv("PORT")
	return
}

func openshift_binding() (hostname string, port string) {
	hostname = os.Getenv("HOST")
	port = os.Getenv("PORT")
	return
}


func main() {
	fmt.Printf("Starting Go Say Hello")

	http.HandleFunc("/", sayhello)

	var hostname = "localhost"
	var port = "3337"

	if os.Getenv("ON_HEROKU") == "1" {
		hostname, port = heroku_binding()
	} else {
		if os.Getenv("PORT") != "" {
			hostname, port = openshift_binding()
		}
	}

	fmt.Printf("HOST (%s)", hostname)
	fmt.Printf("PORT (%s)", port)

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
