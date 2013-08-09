package main

import (
  "fmt"
  "strconv"
  "strings"
  "net/http"
  "html/template"
  "io"
  "io/ioutil"
  "github.com/tarm/goserial"
)

var s io.ReadWriteCloser

func loadSerial() {
  // find the serial port the Arduino is on. It's usually /dev/ttyACM*
  files, err := ioutil.ReadDir("/dev")
  if err != nil {
    panic("Unable to list /dev!")
  }
  port := ""
  for _, f := range files {
    if strings.HasPrefix(f.Name(), "ttyACM") {
      port = f.Name()
      break
    }
  }
  if port == "" {
    panic("Arduino not connected!")
  }

  c := &serial.Config{Name: "/dev/" + port, Baud: 115200}
  s, err = serial.OpenPort(c)
  if err != nil {
    panic("Unable to open serial port!");
  }
  fmt.Printf("Successfully connected to Arduino on /dev/%s.\n", port)
}

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

type Color struct {
  R, G, B int
}

var color Color

func index(w http.ResponseWriter, r *http.Request) {
  if (r.Method == "POST") {
    var err error
    color.R, err = strconv.Atoi(r.FormValue("r"))
    if err != nil {
      color.R = 0
    }
    color.G, err = strconv.Atoi(r.FormValue("g"))
    if err != nil {
      color.G = 0
    }
    color.B, err = strconv.Atoi(r.FormValue("b"))
    if err != nil {
      color.B = 0
    }
    s.Write([]byte{byte(color.R), byte(color.G), byte(color.B), 1})
  }

  err := templates.ExecuteTemplate(w, "index.html", color)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func main() {
  loadSerial()
  loadTemplates()

  fmt.Println("Opening server on localhost:8080")
  http.HandleFunc("/", index)
  http.ListenAndServe(":8080", nil)
}

