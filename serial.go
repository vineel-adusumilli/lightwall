package main

import (
  "fmt"
  "strings"
  "io"
  "io/ioutil"

  "github.com/tarm/goserial"
)

var s io.ReadWriteCloser
var readySerial = make(chan bool)
var lightQueue = make(chan []byte, 100)

func loadSerial() {
  // find the serial port the Arduino is on. It's usually /dev/ttyACM* or /dev/ttyUSB*
  files, err := ioutil.ReadDir("/dev")
  if err != nil {
    panic("Unable to list /dev!")
  }
  port := ""
  for _, f := range files {
    if strings.HasPrefix(f.Name(), "ttyUSB") || strings.HasPrefix(f.Name(), "ttyACM") {
      port = f.Name()
      break
    }
  }
  if port == "" {
    panic("Arduino not connected!")
  }

  // open connection to Arduino
  c := &serial.Config{Name: "/dev/" + port, Baud: 115200}
  s, err = serial.OpenPort(c)
  if err != nil {
    panic("Unable to open serial port!");
  }
  fmt.Printf("Successfully connected to Arduino on /dev/%s.\n", port)

  // start the relevant goroutines
  go readSerial()
  go writeSerial()
}

// monitor for serial readiness
func readSerial() {
  buf := make([]byte, 255)
  readySerial <- true
  for {
    nr, _ := s.Read(buf)
    if nr > 0 {
      readySerial <- true
    }
  }
}

func writeSerial() {
  var color []byte
  for {
    // wait for serial connection to be ready
    <-readySerial

    color = append(<-lightQueue, 0)
    for i, c := range color {
      if c == 1 {
        color[i] = 0
      }
    }
    color[3] = 1

    s.Write(color)
  }
}

