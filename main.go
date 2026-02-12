package main

import (
	"net/http"
)

func mainPage(w http.ResponseWriter, req *http.Request) {

	http.ServeFile(w, req, "index.html")

}

func main() {
	http.Handle("/static/", http.StripPrefix(("/static/"), http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", mainPage)
	http.ListenAndServe(":8090", nil)
}
