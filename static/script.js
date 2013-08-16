window.onload = function() {
  var ready = false;
  var socket = new WebSocket("ws://" + location.host + "/ws");
  var r = document.getElementById("r");
  var g = document.getElementById("g");
  var b = document.getElementById("b");

  socket.onmessage = function(msg) {
    if (msg.data === "ready") {
      ready = true;
      return;
    }

    var rgb = msg.data.split(",");
    r.value = rgb[0];
    g.value = rgb[1];
    b.value = rgb[2];
  };

  socket.onclose = function() {
    if (ready) {
      alert("Connection to the server timed out!");
    } else {
      alert("The server has already accepted the maximum number of connections!");
    }

    ready = false;
  };

  r.onchange = g.onchange = b.onchange = update;

  function update() {
    if (ready) {
      socket.send([r.value, g.value, b.value]);
    }
  }
};
