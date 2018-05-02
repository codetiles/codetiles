const socket = new WebSocket('ws://localhost:8080/api/v1/ws/findgame');

socket.addEventListener('message', function (event) {
  // user isn't logged in
  if (event.data == "User id does not exist") {
    $("#desc").html(`<a href='/'>Back to home page</a>`);
  }
  // game has started, redirect
  if (event.data == "...") {
    window.location.href="/game"
  }
  // if game is starting, hide player count required
  if (event.data.includes('Game starts in')) {
    $('#desc').hide();
  } else {
    $('#desc').show();
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
