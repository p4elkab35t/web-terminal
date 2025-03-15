package main

import (
	"html/template"
	"io"
	"reflect"
	"strings"
	"web-terminal/apps"

	"net/http"
	"path"
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

	words := strings.Split(commandText, " ")

	args := make([]reflect.Value, len(words)-1)

	for i, _ := range args {
		args[i] = reflect.ValueOf(words[i+1])
	}

	method := findService(words[0], &apps.Apps{})

	if method.IsValid() {
		result := method.Call(args)
		w.Write([]byte(result[0].String()))
		return
	} else {
		http.Error(w, "No Command "+words[0]+" found", http.StatusExpectationFailed)
		return
	}
}

func findService(serviceName string, appsRepo *apps.Apps) reflect.Value {
	method := reflect.ValueOf(appsRepo).MethodByName(serviceName)
	return method
}
