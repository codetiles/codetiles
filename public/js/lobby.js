const socket = new WebSocket('ws://localhost:8080/api/v1/ws/findgame');

socket.addEventListener('message', function (event) {
  // user isn't logged in
  if (event.data == "User id does not exist") {
    window.location.href = "/"
  }
  // game has started, redirect
  if (event.data == "...") {
    window.location.href = "/game"
  }
  // show current message as text
  $("#status").text(function(){
    return event.data
  })
});

socket.addEventListener('open', function (event) {
  socket.send(localStorage.getItem('user_id'));
});

// keep connection alive
setInterval(function(){
  socket.send("...")
}, 250);

$(window).on('beforeunload', function(){
    socket.close();
});
