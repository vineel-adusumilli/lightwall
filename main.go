package main

import (
  "fmt"
  "strings"

  "net/http"
  "html/template"
  "io/ioutil"

  "code.google.com/p/go.net/websocket"
)


var templates *template.Template

func loadTemplates() {
  // find all the files in templates/ ending in .html
  files, err := ioutil.ReadDir("templates/")
  if err != nil {
    panic("Unable to list templates/!")
  }
  templateFiles := make([]string, 0)
  for _, f := range files {
    if strings.HasSuffix(f.Name(), ".html") {
      templateFiles = append(templateFiles, "templates/" + f.Name())
    }
  }

  // actually load the templates
  templates = template.Must(template.ParseFiles(templateFiles...))

  for _, t := range templateFiles {
    fmt.Printf("Loaded %s.\n", t)
  }
}

func index(w http.ResponseWriter, r *http.Request) {
  err := templates.ExecuteTemplate(w, "index.html", nil)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func main() {
  loadSerial()
  loadTemplates()
  // set max connections to 10
  initSemaphore(10)
  go h.run()

  fmt.Println("Opening server on localhost:8080")
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
  http.Handle("/ws", websocket.Handler(websocketServer))
  http.HandleFunc("/", index)
  http.ListenAndServe(":8080", nil)
}

