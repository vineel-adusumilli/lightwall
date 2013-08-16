window.onload = function() {
  var socket = new WebSocket("ws://" + location.host + "/ws");
  var r = document.getElementById("r");
  var g = document.getElementById("g");
  var b = document.getElementById("b");

  socket.onmessage = function(msg) {
    var rgb = msg.data.split(",");
    r.value = rgb[0];
    g.value = rgb[1];
    b.value = rgb[2];
  };

  r.onchange = g.onchange = b.onchange = update;

  function update() {
    socket.send([r.value, g.value, b.value]);
  }
};
