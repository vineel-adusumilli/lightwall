package main

import (
  "fmt"
  "strconv"
  "net/http"
  "io"
  "github.com/tarm/goserial"
)

var s io.ReadWriteCloser
var red, green, blue int

func index(w http.ResponseWriter, r *http.Request) {
  if (r.Method == "POST") {
    var err error
    red, err = strconv.Atoi(r.FormValue("r"))
    if err != nil {
      red = 0
    }
    green, err = strconv.Atoi(r.FormValue("g"))
    if err != nil {
      green = 0
    }
    blue, err = strconv.Atoi(r.FormValue("b"))
    if err != nil {
      blue = 0
    }
    s.Write([]byte{byte(red), byte(green), byte(blue), 1})
  }
  fmt.Fprintf(w, "<p>lightwall</p>\n")
  fmt.Fprintf(w, "<form action=\"/\" method=\"post\">")
  fmt.Fprintf(w, "r: <input type=\"text\" name=\"r\" value=\"%v\" />", red)
  fmt.Fprintf(w, "g: <input type=\"text\" name=\"g\" value=\"%v\" />", green)
  fmt.Fprintf(w, "b: <input type=\"text\" name=\"b\" value=\"%v\" />", blue)
  fmt.Fprintf(w, "<input type=\"submit\" name=\"Submit\" />")
  fmt.Fprintf(w, "</form>")
}

func main() {
  c := &serial.Config{Name: "/dev/ttyACM0", Baud: 115200}
  var err error
  s, err = serial.OpenPort(c)
  if err != nil {
    panic("Unable to open serial port!");
  }
  http.HandleFunc("/", index)
  http.ListenAndServe(":8080", nil)
}

