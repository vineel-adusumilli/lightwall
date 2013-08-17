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
      $("#swatch").hide();
      $("#main").text("Connection to the server has timed out! Try reloading.");
    } else {
      $("#swatch").hide();
      $("#main").text("The server has already accepted the maximum number of connections! Reloading in 5 seconds...");
      setTimeout(function() {
        location.reload();
      }, 5000);
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
