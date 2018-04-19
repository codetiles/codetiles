function checkCookie(name) {
  if(Cookies.get(name) == null) {
    return false;
  } else {
    return Cookies.get(name);
  }
}

function handleDisplayNameForm() {
  let submitted_displayname = document.getElementById('displayname').value;

  if(checkCookie('user_id') === false) {
    
  }
}
