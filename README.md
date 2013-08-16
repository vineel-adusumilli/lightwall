lightwall
=========

lightwall is a web app written in go that makes it easy to color the illuminated wall at [TinyFactory](http://tinyfactory.co/).

lightwall uses websockets to connect to a color picker served to the browser. lightwall connects with an Arduino over a USB connection and pushes RGB values as it recieves them from the browser. The end result is an interactive wall that can be controlled from multiple sources at once.

Usage
-----

Running lightwall is as easy as:

```bash
go get github.com/tarm/goserial
go get code.google.com/p/go.net/websocket 
go build
./lightwall
```

