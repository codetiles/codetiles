let isRegistered;
let user_id;
let displayname;

$(function() {
  // check if localStorage isn't supported
  if(!window.localStorage) window.location.href = '/unsupported';

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
  window.location.href = '/game';
}
