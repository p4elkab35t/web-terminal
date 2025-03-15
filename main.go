package main

import (
	// "database/sql"
	// "fmt"
	"html/template"
	"io"

	// "log"
	"net/http"
	// "strconv"
	// "strings"
	"path"
	// _ "github.com/go-sql-driver/mysql"
)

func main() {

	server := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))

	server.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	server.HandleFunc("GET /", renderTerminal)
	server.HandleFunc("POST /command", commandExecute)

	http.ListenAndServe(":80", server)
}

func renderTerminal(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET this url, dont try your trickstery here", http.StatusMethodNotAllowed)
	}

	templateFilePath := path.Join("templates", "terminal.html")
	terminalTemplate, err := template.ParseFiles(templateFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	if err := terminalTemplate.Execute(w, "Terminal"); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
}

func commandExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST this url, dont try your trickstery here", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	commandTextBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	commandText := string(commandTextBytes)

	if commandText == "errtest" {
		http.Error(w, "Error test successfully", http.StatusExpectationFailed)
		return
	}

	answer := "YAY YOU GOT THE COMMAND: " + commandText

	w.Write([]byte(answer))
}

// func findService(serviceName string) {

// }
