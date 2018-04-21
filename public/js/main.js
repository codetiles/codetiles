let isRegistered;

function onLoad() {
  // check if localStorage isn't supported
  if(!window.localStorage){window.location.href = '/unsupported.html';}

  if(localStorage.getItem('user_id') != null) {
    document.getElementById('login_widget').style.display = 'none';
    isRegistered = true;
    auth();
  } else {
    isRegistered = false;
    document.getElementById('login_widget').style.display = 'inline';
  }
}

// register user/re-authenticate user on submission of display name
function handleDisplayNameSubmission() {
  let submitted_displayname = document.getElementById('displayname').value;
  if(!isRegistered) {
    register(submitted_displayname);
  } else if(isRegistered) {
    document.getElementById('login_widget').style.display = 'none';
    auth();
  }
}

// intended for first time use, when displayname is set for the first time
function register(displayname) {
  let data = `{"DisplayName":"`+displayname+`"}`
  $.post("api/v1/createuser", data, function(data) {
    let response = JSON.parse(data);
    localStorage.setItem('user_id', response.Id);
    document.getElementById('login_widget').style.display = 'none';
  })
}

// intended for returning visitors, to authenticate (check if ID is valid)
function auth() {
  $.getJSON("api/v1/verifyuser/"+localStorage.getItem('user_id'), function(data) {
    // do something with response
  })
}
