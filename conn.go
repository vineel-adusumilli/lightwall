package main

// code mostly stolen from http://gary.beagledreams.com/page/go-websocket-chat.html

import (
  "fmt"
  "strings"
  "strconv"
  "time"

  "code.google.com/p/go.net/websocket"
)

type connection struct {
  // The websocket connection.
  ws *websocket.Conn

  // Buffered channel of outbound messages.
  send chan string

  // Time to wait before closing connection after inactivity
  timeout time.Duration
}

var color = make([]byte, 3)

func (c *connection) reader() {
  message := make([]byte, 1024)

  var timer *time.Timer
  if c.timeout > 0 {
    timer = time.NewTimer(c.timeout)
  }

  for {
    if c.timeout > 0 {
      // check if we've timed out
      select {
      case <-timer.C:
        return
      default:
      }
    }

    nr, err := c.ws.Read(message)
    if err != nil {
      break
    }
    if nr > 0 {
      rgb := strings.Split(string(message[:nr]), ",")
      for i := range color {
        c, err := strconv.Atoi(rgb[i])
        if err != nil {
          color[i] = 0
        }
        color[i] = byte(c)
      }

      if c.timeout > 0 {
        timer.Reset(c.timeout)
      }

      h.broadcast <- fmt.Sprintf("%d,%d,%d", color[0], color[1], color[2])
      lightQueue <- color
    }
  }
  c.ws.Close()
}

func (c *connection) writer() {
  for message := range c.send {
    err := websocket.Message.Send(c.ws, message)
    if err != nil {
      break
    }
  }
  c.ws.Close()
}

// semaphore to limit max connections 
var sem chan bool

// initialize the semaphore
func initSemaphore(maxConnections int) {
  sem = make(chan bool, maxConnections)
  for i := 0; i < maxConnections; i++ {
    sem <- true
  }
}

func websocketServer(ws *websocket.Conn) {
  // make sure there's an open connection first, otherwise disconnect
  select {
  case <-sem:
  default:
    return
  }

  // set connection timeout at 60 seconds
  c := &connection{
    send: make(chan string, 256),
    ws: ws,
    timeout: 60 * time.Second,
  }
  h.register <- c
  defer func() {
    h.unregister <- c
    // return one connection to the semaphore
    // but don't block if it's already full for some reason
    select {
    case sem <- true:
    }
  }()
  go c.writer()
  // tell client that server is ready to accept rgb data
  c.send <- "ready";
  c.reader()
}

