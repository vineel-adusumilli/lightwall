window.onload = function() {
  var socket = new WebSocket("ws://" + location.host + "/ws");
  var r = document.getElementById("r");
  var g = document.getElementById("g");
  var b = document.getElementById("b");

  r.onchange = g.onchange = b.onchange = update;

  function update() {
    console.log("Sending " + [r.value, g.value, b.value]);
    socket.send([r.value, g.value, b.value]);
  }
};
