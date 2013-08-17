window.onload = function() {
  var ready = false;
  var socket = new WebSocket("ws://" + location.host + "/ws");

  socket.onmessage = function(msg) {
    var rgb = msg.data.split(",");
    header.style.background = "rgb(" + rgb.join() + ")";
    ready = true;
  };

  socket.onclose = function() {
    if (ready) {
      alert("Connection to the server timed out!");
    } else {
      alert("The server has already accepted the maximum number of connections!");
    }

    ready = false;
  };

  $("#swatch").colorPicker("/static/images/swatch.jpg", update);

  function update(rgb) {
    if (ready) {
      socket.send([rgb.r, rgb.g, rgb.b]);
    }
  }
};
