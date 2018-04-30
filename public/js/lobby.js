const socket = new WebSocket('ws://localhost:8080/api/v1/findgame/ws');

socket.addEventListener('message', function (event) {
  console.log('Message from server ', event.data);
});

socket.addEventListener('open', function (event) {
  socket.send("3nsvyZ9v")
});

$(window).on('beforeunload', function(){
    socket.close();
});
