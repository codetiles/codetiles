let isRegistered;

function onLoad() {
  if(cookie('user_id') != null) {
    let user_id = cookie('user_id');
    document.getElementById('login_widget').style.display = 'none';
    isRegistered = true;
    auth();
  } else {
    isRegistered = false;
    document.getElementById('login_widget').style.display = 'inline';
  }
}

function cookie(name) {
  if(Cookies.get(name) == null) {
    return null;
  } else {
    return Cookies.get(name);
  }
}

// register user/re-authenticate user on submission of display name
function handleDisplayNameSubmission() {
  let submitted_displayname = document.getElementById('displayname').value;
  if(isRegistered) {
    document.getElementById('login_widget').style.display = 'none';
    auth();
  } else if(!isRegistered) {
    register();
  }
}

function register() {
  // intended for first time use, when displayname is set for the first time
}

function auth() {
  // intended for returning visitors, to authenticate (check if ID is valid)
}
