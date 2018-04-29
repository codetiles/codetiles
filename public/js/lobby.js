const socket = new WebSocket('ws://localhost:8080/api/v1/waitforgame/ws')

socket.addEventListener('message', function (event) {
    console.log('Message from server ', event.data);
});
