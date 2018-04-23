let isRegistered;
let user_id;
let displayname;

function onLoad() {
  // check if localStorage isn't supported
  if(!window.localStorage) window.location.href = '/unsupported';

  if(encodeURIComponent(localStorage.getItem('user_id')) != null) {
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
  let submitted_displayname = encodeURIComponent(document.getElementById('displayname').value);
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
    document.getElementById('logged_in_text').innerHTML = ('Logged in as '+displayname+', <a href="#" onclick="logout();">logout</a>');
  })
}

// intended for returning visitors, to authenticate (check if ID is valid)
function auth() {
  $.getJSON("api/v1/verifyuser/"+encodeURIComponent(localStorage.getItem('user_id')), function(data) {
    if(data.Exists === "false") {
      localStorage.removeItem('user_id');
      document.getElementById('login_widget').style.display = 'inline';
    } else if(data.Exists === "true") {
      user_id = encodeURIComponent(localStorage.getItem('user_id'));
      displayname = encodeURIComponent(data.DisplayName);
      document.getElementById('logged_in_text').innerHTML = ('Logged in as '+displayname+', <a href="" onclick="logout();">logout</a>');
    }
  })
}

function logout() {
  localStorage.removeItem('user_id');
  user_id = null;
  displayname = null;
  location.reload();
}
