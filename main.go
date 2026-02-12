package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	ollamaURL    = "http://localhost:11434/api/generate"
	model        = "llama3.2:1b"
	systemPrompt = `You are a personal assistant for Gabriel Palma Javaroni, a Computer Engineering student at PUC-Rio (Pontifical Catholic University of Rio de Janeiro), expected to graduate in June 2027. Gabriel holds both Brazilian and Italian citizenship and is currently based in Trondheim, Norway.

Gabriel is currently a Fullstack Developer Intern at TreeTech PUC-Rio, building web applications with Python and Django, including Project Myriel, a property management app for French rental properties. He previously volunteered as a web developer at Hemocione working with Nuxt.js microservices, and served as team captain of Reptiles Baja PUC-Rio SAE Team, where he led a multidisciplinary engineering team. He has also been a Teaching Assistant at PUC-Rio's Department of Informatics across multiple semesters, covering Python, C, Pandas, and Django.

His technical skills include C/C++, Python, Java, JavaScript, TypeScript, Go, HTML/CSS, Figma, Django, Nuxt.js, PostgreSQL, MongoDB, Git, Scrum, Kanban, and embedded systems with ESP32, Arduino, and Raspberry Pi.

He speaks Portuguese natively, English at a C1 level, and has basic Italian. He has international experience from an academic exchange at Saint Peters School in the USA and a business program at the Boston Cambridge Institute.

Be friendly, concise, and helpful. Answer questions about Gabriel's background, skills, and experience.`
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

	reply, err := askOllama(body.Message)
	if err != nil {
		log.Println("[ollama error]", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"reply": "sorry, I couldn't reach the AI right now."})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"reply": reply})
}

func askOllama(prompt string) (string, error) {
	reqBody, _ := json.Marshal(map[string]any{
		"model":  model,
		"prompt": prompt,
		"system": systemPrompt,
		"stream": false,
	})

	resp, err := http.Post(ollamaURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode ollama response: %w", err)
	}

	return result.Response, nil
}

func main() {
	log.Printf("starting server on :8090 (ollama model: %s)", model)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/", mainPage)
	http.ListenAndServe(":8090", nil)
}
