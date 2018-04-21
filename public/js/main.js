let isRegistered;

function onLoad() {
  let parsedCookie = {};
  for(let i = 0;;i++){
    let part = document.cookie.split("; ")[i].split("=");
    if (part == "") break;
    parsedCookie[part[0]] = part[1];
  }

  if(parsedCookie.user_id != null) {
    let user_id = parsedCookie.user_id;
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
  if(isRegistered) {
    document.getElementById('login_widget').style.display = 'none';
    auth();
  } else if(!isRegistered) {
    register(submitted_displayname);
  }
}

// intended for first time use, when displayname is set for the first time
function register(displayname) {
  let data = `{"DisplayName":"`+displayname+`"}`
  $.post("api/v1/createuser", data, function(data) {
    console.log(data);
  })
}

// intended for returning visitors, to authenticate (check if ID is valid)
function auth() {
  $.getJSON("api/v1/verifyuser/"+user_id, function(data) {
    // do something with response
  })
}
