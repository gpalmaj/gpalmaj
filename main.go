package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func mainPage(w http.ResponseWriter, req *http.Request) {

	http.ServeFile(w, req, "index.html")

}

func chatHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	log.Println("[chat]", body.Message)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"reply": "thanks for reaching out! I'll get back to you soon."})
}

func main() {
	http.Handle("/static/", http.StripPrefix(("/static/"), http.FileServer(http.Dir("static"))))
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/", mainPage)
	http.ListenAndServe(":8090", nil)
}
