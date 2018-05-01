function checkLogin() {
  // if user id isn't defined, redirect to homepage
  if (localStorage.getItem('user_id') === null) {
    window.location.href = '/';
  } else {
    // if user id is defined, check with server to verify it exists
    $.getJSON("/api/v1/verifyuser/"+localStorage.getItem('user_id'), function(data) {
      if(data.Exists == "false") {
        localStorage.removeItem('user_id');
        window.location.href = '/';
      }
    });
  }
}
