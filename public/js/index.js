let isRegistered;
let user_id;
let displayname;

$(function() {
  // check if localStorage isn't supported
  if(!window.localStorage) window.location.href = '/unsupported';

  // if enter key is pressed (while displayname is focused), submit data
  $("#displayname").on('keyup', function (e) {
    if (e.keyCode == 13) {
      handleDisplayNameSubmission();
    }
  });

  if(localStorage.getItem('user_id') != null) {
    $('#login_widget').hide();
    isRegistered = true;
    auth();
  } else {
    isRegistered = false;
    $('#login_widget').show();
  }
});

// register user/re-authenticate user on submission of display name
function handleDisplayNameSubmission() {
  let submitted_displayname = $('#displayname').val();
  if(!isRegistered) {
    if(submitted_displayname !== '' && /\s/.test(submitted_displayname) !== true) {
      register(submitted_displayname);
    } else {
      $('#login_widget').append('<div class="error">Name cannot be blank or contain spaces.</div>');
    }
  } else if(isRegistered) {
    $('#login_widget').hide();
    auth();
  }
}

// intended for first time use, when displayname is set for the first time
function register(displayname) {
  let data = `{"DisplayName":"`+displayname+`"}`
  $.post("api/v1/createuser", data, function(data) {
    let response = JSON.parse(data);
    localStorage.setItem('user_id', response.Id);
    $('#login_widget').hide();
    $("#logged_in_text").html(`Logged in as `+displayname+`, <a href="" onclick="logout();">logout</a><br><button type="submit" onClick="startGame()">Start a Game</button>`);
  })
}

// intended for returning visitors, to authenticate (check if ID is valid)
function auth() {
  $.getJSON("api/v1/verifyuser/"+localStorage.getItem('user_id'), function(data) {
    if(data.Exists == "false") {
      localStorage.removeItem('user_id');
      $('#login_widget').show();
      location.reload();
    } else {
      user_id = localStorage.getItem('user_id');
      displayname = data.DisplayName;
      $("#logged_in_text").html(`Logged in as `+displayname+`, <a href="" onclick="logout();">logout</a><br><button type="submit" onClick="startGame()">Start a Game</button>`);
    }
  })
}

function logout() {
  localStorage.removeItem('user_id');
  user_id = null;
  displayname = null;
  location.reload();
}

function startGame() {
  $("#code_textarea").text("")
  $("#code_textarea").animate({
    "width": "0vw",
    "padding" : 0,
    "margin" : 0
  }, 300);
  $("#code").animate({
    "width": "100vw",
  });

  $(".slogan").text("Game will start when 2 players are in queue.");

  const socket = new WebSocket('ws://localhost:8080/api/v1/ws/findgame');

  socket.addEventListener('message', function (event) {
    // user isn't logged in
    if (event.data == "User id does not exist") {
      $(".slogan").html(`<a href='/'>Back to home page</a>`);
    }
    // game has started, redirect
    if (event.data == "...") {
      window.location.href="/game"
    }
    // if game is starting, hide player count required
    if (event.data.includes('Game starts in')) {
      $('.slogan').hide();
    } else {
      $('.slogan').show();
    }
    // show current message as text
    $(".title").text(event.data);
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

  $("#logged_in_text").hide();
  $("#login_widget").hide();
}
