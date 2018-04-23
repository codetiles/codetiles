let isRegistered;
let user_id;
let displayname;

onLoad();

function onLoad() {
  console.log('onload ran')
  // check if localStorage isn't supported
  if(!window.localStorage) window.location.href = '/unsupported';

  if(localStorage.getItem('user_id') != null) {
    console.log('user id exists')
    $('#login_widget').hide();
    isRegistered = true;
    auth();
  } else {
    console.log('user id does not exist')
    isRegistered = false;
    $('#login_widget').show();
  }
}

// register user/re-authenticate user on submission of display name
function handleDisplayNameSubmission() {
  console.log('button pressed')
  let submitted_displayname = $('#displayname').text();
  console.log('displayname = '+displayname)
  if(!isRegistered) {
    console.log('not registered, registering..')
    register(submitted_displayname);
  } else if(isRegistered) {
    console.log('registered, hiding login & authenticating..')
    $('#login_widget').hide();
    auth();
  }
}

// intended for first time use, when displayname is set for the first time
function register(displayname) {
  console.log('registering..')
  let data = `{"DisplayName":"`+displayname+`"}`
  console.log('making a request..')
  $.post("api/v1/createuser", data, function(data) {
    console.log('request done')
    let response = JSON.parse(data);
    console.log('response: '+response)
    localStorage.setItem('user_id', response.Id);
    $('#login_widget').hide();
    $("#logged_in_text").text(`Logged in as `+displayname+`, <a href="" onclick="logout();">logout</a>`);
  })
}

// intended for returning visitors, to authenticate (check if ID is valid)
function auth() {
  console.log('authenticating..')
  $.getJSON("api/v1/verifyuser/"+localStorage.getItem('user_id'), function(data) {
    console.log('request done, exists: '+data.Exists)
    if(data.Exists == "false") {
      console.log('does not exist, removing user & showing login')
      localStorage.removeItem('user_id');
      $('#login_widget').show();
    } else {
      console.log('does exist, storing userid + displayname & showing as logged in')
      user_id = localStorage.getItem('user_id');
      displayname = encodeURIComponent(data.DisplayName);
      $("#logged_in_text").text(`Logged in as `+displayname+`, <a href="" onclick="logout();">logout</a>`);
    }
  })
}

function logout() {
  console.log('logging out.. removing localstorage item and refreshing page')
  localStorage.removeItem('user_id');
  user_id = null;
  displayname = null;
  location.reload();
}
