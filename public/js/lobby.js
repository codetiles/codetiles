const socket = new WebSocket('ws://localhost:8080/api/v1/ws/findgame');

socket.addEventListener('message', function (event) {
  if (event.data == "User id does not exist") {
    window.location.href = "/"
  }
  $("#status").text(function(){
    return event.data
  })
});

socket.addEventListener('open', function (event) {
  socket.send(localStorage.getItem('user_id'));
});

setInterval(function(){
  socket.send("...")
}, 250);

$(window).on('beforeunload', function(){
    socket.close();
});
